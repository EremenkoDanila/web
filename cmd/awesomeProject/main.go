package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "lab1/internal/app/config"
    "lab1/internal/app/dsn"
    "lab1/internal/app/handler"
    "lab1/internal/app/repository"
    "lab1/internal/pkg"
    "lab1/internal/app/s3"
    "lab1/internal/app/redis"
    "context"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
    _ "lab1/docs" // <- обязательно для swag
)

// @title Lab4 API
// @version 1.0
// @description API for Lab1 project
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
    router := gin.Default()

    url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // ссылка на Swagger JSON
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))


    conf, err := config.NewConfig()
    if err != nil {
        logrus.Fatalf("error loading config: %v", err)
    }

    postgresString := dsn.FromEnv()
    fmt.Println("Database DSN:", postgresString)

    rep, err := repository.New(postgresString)
    if err != nil {
        logrus.Fatalf("error initializing repository: %v", err)
    }

    minioClient, err := s3.InitMinioClient(&conf.MinioConfig)
    if err != nil {
        logrus.Fatalf("error initializing MinIO: %v", err)
    }


    ctx := context.Background()
    redisClient, err := redis.New(ctx, conf.Redis)
    if err != nil {
        logrus.Fatalf("error initializing Redis: %v", err)
    }
    defer func() {
        if err := redisClient.Close(); err != nil {
            logrus.Errorf("error closing Redis client: %v", err)
        }
    }()


    hand := handler.NewHandler(rep, minioClient, &conf.MinioConfig, conf, redisClient)
    app := pkg.NewApp(conf, router, hand)
    app.RunApp()
}
