package auth

import "github.com/gin-gonic/gin"

func TokenAuth(optional bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
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
