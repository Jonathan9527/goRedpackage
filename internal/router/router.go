package router

import (
	"learnGO/internal/handler"
	"learnGO/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	Home       *handler.HomeHandler
	User       *handler.UserHandler
	RedPackage *handler.RedPackageHandler
}

func Register(engine *gin.Engine, handlers Handlers) {
	engine.GET("/", handlers.Home.Greeting)
	engine.GET("/health", handlers.Home.Health)

	engine.GET("/users", handlers.User.List)
	engine.GET("/users/:account", UserDetailRateLimiter(10, 5), handlers.User.FindByAccount)
	v1 := engine.Group("/v1")
	{
		// 用户相关接口
		redPackage := v1.Group("/redpackage")
		redPackage.POST("/sendRedPackage", handlers.RedPackage.SendRedPackage)
		redPackage.GET("/getRedPackage", handlers.RedPackage.GetRedPackage)

	}
}

func UserDetailRateLimiter(capacity int, refillRate int) gin.HandlerFunc {
	return middleware.TokenBucketRateLimiter(middleware.NewTokenBucket(capacity, float64(refillRate)))
}
