package auth

import (
	"errors"
	"os"
	"phrasmotica/bore-score-api/models"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type JWTClaim struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

const tokenLifetime = 1 * time.Hour

func GenerateJWT(user *models.User) (tokenString string, err error) {
	expirationTime := time.Now().Add(tokenLifetime)

	claims := &JWTClaim{
		Email:       user.Email,
		Username:    user.Username,
		Permissions: user.Permissions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

// https://www.sohamkamani.com/golang/jwt-authentication/
func RefreshJWT(tokenStr string) (newToken string, err error) {
	claims := &JWTClaim{}

	newToken = ""

	parsedToken, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	if !parsedToken.Valid {
		return
	}

	// don't refresh if more than 30 seconds until expiry
	if time.Until(time.Unix(claims.ExpiresAt, 0)) > 30*time.Second {
		newToken = tokenStr
		return
	}

	expirationTime := time.Now().Add(tokenLifetime)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	newToken, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error, claims *JWTClaim) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	return
}
