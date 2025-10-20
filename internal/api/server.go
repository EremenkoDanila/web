package api

import (
	"html/template"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"lab1/internal/app/handler"
	"lab1/internal/app/repository"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория: ", err)
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"lower": strings.ToLower,
	})

	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/it_soft", handler.GetOrders)
	r.GET("/Software/:id", handler.GetOrderByID)
	r.GET("/Software_request/:count/:id", handler.GetCartItems)


	r.Run(":8080")
	log.Println("Server down")
}
