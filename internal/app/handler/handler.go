package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    "lab1/internal/app/config"
    "lab1/internal/app/repository"
)

type Handler struct {
    Repository  *repository.Repository
    MinioClient *minio.Client
    MinioConfig *config.MinioConfig
}

func NewHandler(r *repository.Repository, minioClient *minio.Client, minioConfig *config.MinioConfig) *Handler {
    return &Handler{
        Repository:  r,
        MinioClient: minioClient,
        MinioConfig: minioConfig,
    }
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
    router.GET("/api/software", h.GetSoftwareWithFilter)
    router.GET("/api/software/:id", h.GetSoftwareByID)
    router.POST("/api/software", h.AddSoftware)
    router.POST("/api/software/:id/image", h.AddPicture)
    router.PUT("/api/software/:id", h.ChangeSoftware)
    router.DELETE("/api/software/:id", h.DeleteSoftware)
	router.POST("/api/software/to_request/draft/:software_id", h.AddToDraftRequest)


    router.GET("/api/request/cart", h.GetUserCart)           
    router.GET("/api/requests", h.GetRequestsList)          
    router.GET("/api/request/:id", h.GetRequestByID)     
    router.PUT("/api/requests/:id/fields", h.UpdateRequestFields)
    router.PUT("/api/requests/:id/form", h.FormRequest)
    router.DELETE("/api/requests/:id", h.DeleteRequest)
    router.PUT("/api/requests/:id/status/:status", h.UpdateRequestStatusByModerator)



    router.DELETE("/api/InstallationTime/:requests_id/software/:software_id", h.DeleteInstallationTime)
    router.PUT("/api/InstallationTime/:request_id/software/:software_id/install_time", h.UpdateInstallationTime)


    router.POST("/api/user/add", h.RegisterUser)
    router.GET("/api/user", h.GetUserProfile)      // получить данные пользователя
    router.PUT("/api/user", h.UpdateUserProfile)   // обновить данные пользователя
    router.POST("/api/user/login", h.Login)        // аутентификация
    router.POST("/api/user/logout", h.Logout)      // деавторизация
}