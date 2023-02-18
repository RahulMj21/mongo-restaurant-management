package middlewares

import (
	"github.com/gin-gonic/gin"
)

func Authentication(c *gin.Context) {
	c.Next()
}
