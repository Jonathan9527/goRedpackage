package main

import (
	"log"

	"learnGO/internal/config"
	"learnGO/internal/database"
	"learnGO/internal/handler"
	"learnGO/internal/repository"
	approuter "learnGO/internal/router"
	"learnGO/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql db: %v", err)
	}
	defer sqlDB.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}
	rabbitMQConn, err := database.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf("connect rabbitmq: %v", err)
	}
	defer rabbitMQConn.Close()

	redisClient, err := database.NewRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("connect redis: %v", err)
	}
	defer redisClient.Close()
	router := gin.Default()

	homeService := service.NewHomeService(db)
	homeHandler := handler.NewHomeHandler(homeService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	redPackageService := service.NewRedPackageService(userRepository, rabbitMQConn, redisClient)

	getRedPackageService := service.NewGetRedPackageService(redisClient)
	redPackageHandler := handler.NewRedPackageHandler(redPackageService, getRedPackageService)

	approuter.Register(router, approuter.Handlers{
		Home:       homeHandler,
		User:       userHandler,
		RedPackage: redPackageHandler,
	})

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
