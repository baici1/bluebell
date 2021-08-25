package router

import (
	"bluebell/controllers"
	"bluebell/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.POST("/signup", controllers.SignUpHandler)
	r.POST("/login", controllers.LoginHandler)
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
