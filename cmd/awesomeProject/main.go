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
    "lab1/internal/s3"
)

func main() {
    router := gin.Default()

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

    hand := handler.NewHandler(rep, minioClient, &conf.MinioConfig)
    app := pkg.NewApp(conf, router, hand)
    app.RunApp()
}
