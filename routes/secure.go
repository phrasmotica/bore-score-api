package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"message": "pong", "username": c.GetString("username")})
}
