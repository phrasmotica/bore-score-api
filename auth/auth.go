package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("JWT_PUBLIC_KEY"))

type JWTClaim struct {
	Username       string                          `json:"preferred_username"`
	Email          string                          `json:"email"`
	Azp            string                          `json:"azp"`
	ResourceAccess map[string]ClientResourceAccess `json:"resource_access"`
	jwt.StandardClaims
}

type ClientResourceAccess struct {
	Roles []string `json:"roles"`
}

func ValidateToken(signedToken string) (claims *JWTClaim, err error) {
	publicKey, err := parseKeycloakRSAPublicKey(string(jwtKey))
	if err != nil {
		panic(err)
	}

	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// return the public key that is used to validate the token.
			return publicKey, nil
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

func parseKeycloakRSAPublicKey(base64Encoded string) (*rsa.PublicKey, error) {
	buf, err := base64.StdEncoding.DecodeString(base64Encoded)
	if err != nil {
		return nil, err
	}

	parsedKey, err := x509.ParsePKIXPublicKey(buf)
	if err != nil {
		return nil, err
	}

	publicKey, ok := parsedKey.(*rsa.PublicKey)
	if ok {
		return publicKey, nil
	}

	return nil, fmt.Errorf("unexpected key type %T", publicKey)
}
