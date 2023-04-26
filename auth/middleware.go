package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// taken from https://stackoverflow.com/a/29439630
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func TokenAuth(optional bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		success, tokenString := parseToken(c)

		if !success {
			if optional {
				c.Next()
				return
			}

			c.IndentedJSON(401, gin.H{"error": "request does not contain an access token"})
			c.Abort()
			return
		}

		err, claims := ValidateToken(tokenString)
		if err != nil {
			c.IndentedJSON(401, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}

func parseToken(c *gin.Context) (bool, string) {
	header := c.GetHeader("Authorization")
	splitToken := strings.Split(header, "Bearer ")

	if len(splitToken) != 2 {
		return false, ""
	}

	return true, splitToken[1]
}
