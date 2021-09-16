package router

import (
	"bluebell/controllers"
	"bluebell/logger"
	"bluebell/middlewares"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func SetUp() *gin.Engine {
	r := gin.New()
	//每两秒限制一个
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(time.Second*2, 1))
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controllers.SignUpHandler)
	v1.POST("/login", controllers.LoginHandler)
	//jwt中间件
	v1.Use(middlewares.JWTAuthMiddleware())

	{
		v1.GET("/community", controllers.CommunityHandler)
		v1.GET("/community/:id", controllers.CommunityDetailHandler)
		v1.POST("/post", controllers.CreatePostHandler)
		v1.GET("/post/:id", controllers.GetPostDetailHandler)
		v1.GET("/posts", controllers.GetPostListHandler)
		//根据时间或分数获取帖子列表
		v1.GET("/posts2", controllers.GetPostListHandler2)
		v1.POST("/vote", controllers.PostVoteController)
	}
	pprof.Register(r) //注册pprof路由
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
