package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/minio/minio-go/v7"
    "lab1/internal/app/config"
    "lab1/internal/app/repository"
	"lab1/internal/app/middlewares"
	"lab1/internal/app/role"
    "lab1/internal/app/redis"
)

type Handler struct {
    Repository  *repository.Repository
    MinioClient *minio.Client
    MinioConfig *config.MinioConfig
    Config      *config.Config
    Redis       *redis.Client
}

func NewHandler(r *repository.Repository, minioClient *minio.Client, minioConfig *config.MinioConfig, cfg *config.Config,  redisClient *redis.Client) *Handler {
	return &Handler{
		Repository:  r,
		MinioClient: minioClient,
		MinioConfig: minioConfig,
		Config:      cfg,
        Redis:       redisClient,
	}
}


func (h *Handler) RegisterHandler(router *gin.Engine) {
    // --- Гость (не авторизован) ---
    guest := router.Group("/api")
    {
        guest.POST("/user/login", h.Login)
        guest.POST("/user/add", h.RegisterUser)
        guest.GET("/software", h.GetSoftwareWithFilter)
        guest.GET("/software/:id", h.GetSoftwareByID)
    }

    // --- Creater (вошедший пользователь) ---
    creater := router.Group("/api")
    creater.Use(middlewares.WithAuthCheck(h.Config.JWT.Token, h.Redis, role.Creater, role.Moderator))
    {
        creater.GET("/user", h.GetUserProfile) // как? 
        creater.PUT("/user", h.UpdateUserProfile) //+

        creater.GET("/request/cart", h.GetUserCart) //+
        creater.GET("/requests", h.GetRequestsList) //+
        creater.GET("/request/:id", h.GetRequestByID) // ??


        creater.PUT("/requests/:id/fields", h.UpdateRequestFields)  // ?
        creater.PUT("/requests/form", h.FormRequest)  //?
        creater.DELETE("/requests/:id", h.DeleteRequest)  //?


        creater.DELETE("/InstallationTime/:requests_id/software/:software_id", h.DeleteInstallationTime)  // ?
        creater.PUT("/InstallationTime/:request_id/software/:software_id/install_time", h.UpdateInstallationTime) // ?
        
    
        creater.POST("/user/logout", h.Logout)
        creater.POST("/software/to_request/draft/:software_id", h.AddToDraftRequest) 
    }

    // --- Moderator ---
    moderator := router.Group("/api")
    moderator.Use(middlewares.WithAuthCheck(h.Config.JWT.Token, h.Redis, role.Moderator))
    {
        moderator.POST("/software", h.AddSoftware)
        moderator.POST("/software/:id/image", h.AddPicture)
        moderator.PUT("/software/:id", h.ChangeSoftware)
        moderator.DELETE("/software/:id", h.DeleteSoftware)
        moderator.PUT("/requests/:id/status/:status", h.UpdateRequestStatusByModerator)


        
    }
}
