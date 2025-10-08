package api

import (
	"html/template"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"log"
	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
	"strings"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	
	// Добавляем функцию lower в шаблоны с правильным типом
	r.SetFuncMap(template.FuncMap{
		"lower": strings.ToLower,
	})
	
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/it_soft", handler.GetOrders)
	r.GET("/apps/:title", handler.GetOrder)
	r.GET("/cart", handler.GetCartItems)

	r.Run()
	log.Println("Server down")
}