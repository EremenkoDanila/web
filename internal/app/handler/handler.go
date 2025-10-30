package handler

import (
	"github.com/gin-gonic/gin"
	"lab1/internal/app/repository"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/software", h.GetAllSoftware)
	router.GET("/software/:id", h.GetSoftwareByID)
	router.GET("/software_request/:id", h.GetCartByRequestID)
	router.POST("/software/add_to_request", h.AddSoftwareToRequest)
	router.POST("/software_request/delete", h.DeleteRequest)
	router.POST("/software_request/update_time", h.UpdateInstallTime)
	router.POST("/software_request/update_phone", h.UpdatePhone)
}


func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
	router.Static("/img", "./resources/img")
}
