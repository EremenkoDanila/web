package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"lab1/internal/app/ds"
	"lab1/internal/app/role"
	"lab1/internal/app/s3"
	auth "lab1/internal/app/usercontext"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	// "crypto/sha1"
	// "encoding/hex"
	_ "lab1/internal/app/dto"
	"log"
	"strings"
)

// ====================== SOFTWARE ======================

// GetSoftwareWithFilter godoc
// @Summary Get list of software
// @Description Get list of software with optional filtering by title
// @Tags Software
// @Accept json
// @Produce json
// @Param app query string false "Software title filter"
// @Success 200 {array} dto.SoftwareDTO
// @Failure 500 {object} map[string]string
// @Router /api/software [get]
func (h *Handler) GetSoftwareWithFilter(c *gin.Context) {
	title := c.Query("app")

	softwares, err := h.Repository.GetSoftwaresWithFilter(title)
	if err != nil {
		logrus.Errorf("Failed to get softwares: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get softwares"})
		return
	}

	c.JSON(http.StatusOK, softwares)
}

// GetSoftwareByID godoc
// @Summary Get software by ID
// @Description Get software by its ID
// @Tags Software
// @Accept json
// @Produce json
// @Param id path int true "Software ID"
// @Success 200 {object} dto.SoftwareDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/software/{id} [get]
func (h *Handler) GetSoftwareByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	software, err := h.Repository.GetSoftwareByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	c.JSON(http.StatusOK, software)
}

// AddSoftware godoc
// @Summary Add new software
// @Description Create new software entry
// @Tags Software
// @Accept json
// @Produce json
// @Param software body dto.SoftwareDTO true "Software info"
// @Success 201 {object} dto.SoftwareDTO
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/software [post]
func (h *Handler) AddSoftware(c *gin.Context) {
	var software ds.Software

	if err := c.BindJSON(&software); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if software.Title == "" || software.OS == "" || software.Version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title, OS and Version are required"})
		return
	}

	if err := h.Repository.AddSoftware(&software); err != nil {
		logrus.Errorf("Failed to add software: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add software"})
		return
	}

	c.JSON(http.StatusCreated, software)
}

// AddPicture godoc
// @Summary Upload software picture
// @Description Upload image for software
// @Tags Software
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Software ID"
// @Param image formData file true "Image file"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/software/{id}/image [post]
func (h *Handler) AddPicture(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	imgURL, size, filename, err := s3.UploadSoftwareImage(
		c.Request.Context(),
		h.Repository,
		h.MinioClient,
		h.MinioConfig,
		uint(id),
		file,
	)
	if err != nil {
		logrus.Errorf("Upload error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Image uploaded successfully",
		"image": gin.H{
			"url":      imgURL,
			"filename": filename,
			"size":     size,
		},
	})
}

// ChangeSoftware godoc
// @Summary Update software
// @Description Update existing software by ID
// @Tags Software
// @Accept json
// @Produce json
// @Param id path int true "Software ID"
// @Param software body dto.SoftwareDTO true "Updated software info"
// @Success 200 {object} dto.SoftwareDTO
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/software/{id} [put]
func (h *Handler) ChangeSoftware(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	existingSoftware, err := h.Repository.GetSoftwareByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	var updateData ds.Software
	if err := c.BindJSON(&updateData); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if updateData.Title != "" {
		existingSoftware.Title = updateData.Title
	}
	if updateData.Description != "" {
		existingSoftware.Description = updateData.Description
	}
	if updateData.FullDescription != "" {
		existingSoftware.FullDescription = updateData.FullDescription
	}
	if updateData.OS != "" {
		existingSoftware.OS = updateData.OS
	}
	if updateData.Size != 0 {
		existingSoftware.Size = updateData.Size
	}
	if updateData.Version != "" {
		existingSoftware.Version = updateData.Version
	}
	if updateData.Functions != "" {
		existingSoftware.Functions = updateData.Functions
	}

	if err := h.Repository.UpdateSoftware(existingSoftware); err != nil {
		logrus.Errorf("Failed to update software: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update software"})
		return
	}

	c.JSON(http.StatusOK, existingSoftware)
}

// DeleteSoftware godoc
// @Summary Soft delete software
// @Description Soft delete software by ID
// @Tags Software
// @Accept json
// @Produce json
// @Param id path int true "Software ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/software/{id} [delete]
func (h *Handler) DeleteSoftware(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Получаем программу и её картинку
	software, err := h.Repository.GetSoftwareByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	// Удаляем изображение из MinIO (если есть)
	if software.ImgURL != "" {
		if err := s3.DeleteSoftwareImage(c.Request.Context(), h.MinioClient, h.MinioConfig, software.ImgURL); err != nil {
			logrus.Warnf("Failed to delete image from MinIO: %v", err)
			// продолжаем, даже если картинка не удалилась
		}
	}

	// Мягкое удаление из базы
	if err := h.Repository.SoftDeleteSoftware(uint(id)); err != nil {
		logrus.Errorf("Failed to delete software: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete software"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Software deleted successfully"})
}

// AddToDraftRequest godoc
// @Summary Add software to user's draft request
// @Description Add a software item to the current user's draft request (cart). Creates draft if none exists.
// @Tags Software
// @Accept json
// @Produce json
// @Param software_id path int true "Software ID"
// @Success 200 {object} map[string]interface{} "Added software and updated draft request"
// @Failure 400 {object} map[string]string "Invalid software ID"
// @Failure 404 {object} map[string]string "Software not found"
// @Failure 401 {object} map[string]string "Unauthorized: user UUID missing or invalid"
// @Failure 403 {object} map[string]string "Forbidden: invalid or missing Authorization header"
// @Failure 500 {object} map[string]string "Failed to create draft request or add software"
// @Security BearerAuth
// @Router /api/software/to_request/draft/{software_id} [POST]
func (h *Handler) AddToDraftRequest(c *gin.Context) {
	// 1️⃣ Получаем UUID пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем пользователя из базы по UUID
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}
	userID := currentUser.UserID

	// 3️⃣ Получаем ID софта из параметра
	softwareIDStr := c.Param("software_id")
	softwareID, err := strconv.ParseUint(softwareIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid software ID"})
		return
	}

	// 4️⃣ Проверяем существование софта
	software, err := h.Repository.GetSoftwareByID(uint(softwareID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
		return
	}

	// 5️⃣ Получаем или создаём черновик заявки
	draftRequest, err := h.Repository.GetDraftRequestByUser(userID)
	if err != nil {
		draftRequest, err = h.Repository.CreateDraftRequest(userID)
		if err != nil {
			logrus.Errorf("Failed to create draft request: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create draft request"})
			return
		}
		logrus.Infof("Created new draft request for user %d", userID)
	}

	// 6️⃣ Добавляем софт в заявку
	err = h.Repository.AddSoftwareToRequest(draftRequest.RequestID, uint(softwareID))
	if err != nil {
		logrus.Errorf("Failed to add software to request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add software to draft request"})
		return
	}

	// 7️⃣ Получаем установочные времена
	installationTimes, err := h.Repository.GetInstallationTimesWithSoftware(draftRequest.RequestID)
	if err != nil {
		logrus.Errorf("Failed to get installation times: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve installation details"})
		return
	}

	// 8️⃣ Формируем ответ
	var phone string
	if draftRequest.Phone.Valid {
		phone = draftRequest.Phone.String
	}

	finishDt := ""
	if draftRequest.FinishDt.Valid {
		finishDt = draftRequest.FinishDt.Time.Format(time.RFC3339)
	}

	var softwares []gin.H
	for _, inst := range installationTimes {
		s := inst.Software

		var installTimeStr, finalTimeStr string
		if inst.InstallTime != nil {
			installTimeStr = inst.InstallTime.Format(time.RFC3339)
		}
		if inst.FinalTime != nil {
			finalTimeStr = inst.FinalTime.Format(time.RFC3339)
		}

		softwares = append(softwares, gin.H{
			"software_id":  s.SoftwareId,
			"title":        s.Title,
			"img_url":      s.ImgURL,
			"os":           s.OS,
			"version":      s.Version,
			"install_time": installTimeStr,
			"final_time":   finalTimeStr,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Software added to draft request successfully",
		"total_softwares": len(softwares),
		"added_software": gin.H{
			"software_id": software.SoftwareId,
			"title":       software.Title,
			"img_url":     software.ImgURL,
			"os":          software.OS,
			"version":     software.Version,
		},
		"request": gin.H{
			"request_id": draftRequest.RequestID,
			"status":     draftRequest.Status,
			"create_dt":  draftRequest.CreateDt,
			"update_dt":  draftRequest.UpdateDt,
			"finish_dt":  finishDt,
			"phone":      phone,
			"softwares":  softwares,
		},
	})
}

// GetUserCart godoc
// @Summary Get user's draft request (cart)
// @Description Get the current draft request of the logged-in user, including software count
// @Tags Requests
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Draft request info"
// @Failure 401 {object} map[string]string "User UUID missing in context"
// @Failure 500 {object} map[string]string "Cannot load current user or software count"
// @Security BearerAuth
// @Router /api/requests/cart [get]
func (h *Handler) GetUserCart(c *gin.Context) {
	// 1️⃣ Получаем UUID пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем userID из базы
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}
	userID := currentUser.UserID

	// 3️⃣ Получаем черновую заявку пользователя
	request, err := h.Repository.GetDraftRequestByUser(userID)
	if err != nil {
		// Если черновика нет, возвращаем пустую корзину
		c.JSON(http.StatusOK, gin.H{"request_id": nil, "software_count": 0})
		return
	}

	// 4️⃣ Получаем список software в черновике
	installationTimes, err := h.Repository.GetInstallationTimesWithSoftware(request.RequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get software count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id":     request.RequestID,
		"software_count": len(installationTimes),
	})
}

// GetRequestsList godoc
// @Summary Get list of requests
// @Description Get list of requests, optionally filtered
// @Tags Requests
// @Accept json
// @Produce json
// @Param status query string false "Request status"
// @Param start_date query string false "Start date YYYY-MM-DD"
// @Param end_date query string false "End date YYYY-MM-DD"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests [get]
func (h *Handler) GetRequestsList(c *gin.Context) {
    status := c.Query("status")
    startDateStr := c.Query("start_date")
    endDateStr := c.Query("end_date")

    var startDate, endDate *time.Time
    if startDateStr != "" {
        t, err := time.Parse("2006-01-02", startDateStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
            return
        }
        startDate = &t
    }
    if endDateStr != "" {
        t, err := time.Parse("2006-01-02", endDateStr)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
            return
        }
        endDate = &t
    }

    roleVal, exists := c.Get("role")
    if !exists {
        c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
        return
    }
    userUUIDVal, exists := c.Get("user_uuid")
    if !exists {
        c.JSON(http.StatusForbidden, gin.H{"error": "user_uuid not found"})
        return
    }

    userRole := roleVal.(role.Role)
    userUUID := userUUIDVal.(string)

    var requests []ds.SoftwareRequest
    var err error

    if userRole == role.Moderator {
        requests, err = h.Repository.GetRequestsFiltered(status, startDate, endDate)
    } else {
        requests, err = h.Repository.GetRequestsFilteredByUser(userUUID, status, startDate, endDate)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
        return
    }

    // ---------- DTO ----------
    type RequestResp struct {
        RequestID      uint   `json:"request_id"`
        Status         string `json:"status"`
        CreateDt       string `json:"create_dt"`
        UpdateDt       string `json:"update_dt"`
        FinishDt       string `json:"finish_dt"`
        CreatorID      uint   `json:"creator_id"`
        ModeratorID    uint   `json:"moderator_id"`
        Phone          string `json:"phone"`
        CompletedCount int64  `json:"completed_count"`
    }

    var response []RequestResp

    for _, r := range requests {
        phone := ""
        if r.Phone.Valid {
            phone = r.Phone.String
        }

        finishDt := ""
        if r.FinishDt.Valid {
            finishDt = r.FinishDt.Time.Format(time.RFC3339)
        }

        // <- Вернули!
        completedCount, _ := h.Repository.CountCompletedSoftware(r.RequestID)

        response = append(response, RequestResp{
            RequestID:      r.RequestID,
            Status:         r.Status,
            CreateDt:       r.CreateDt.Format(time.RFC3339),
            UpdateDt:       r.UpdateDt.Format(time.RFC3339),
            FinishDt:       finishDt,
            CreatorID:      r.CreatorID,
            ModeratorID:    r.ModeratorID,
            Phone:          phone,
            CompletedCount: completedCount,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "requests": response,
    })
}


// func (h *Handler) GetRequestsList(c *gin.Context) {
// 	status := c.Query("status")
// 	startDateStr := c.Query("start_date")
// 	endDateStr := c.Query("end_date")

// 	var startDate, endDate *time.Time
// 	if startDateStr != "" {
// 		t, err := time.Parse("2006-01-02", startDateStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
// 			return
// 		}
// 		startDate = &t
// 	}
// 	if endDateStr != "" {
// 		t, err := time.Parse("2006-01-02", endDateStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
// 			return
// 		}
// 		endDate = &t
// 	}

// 	// Получаем роль и user_uuid из контекста
// 	roleVal, exists := c.Get("role")
// 	if !exists {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "role not found"})
// 		return
// 	}
// 	userUUIDVal, exists := c.Get("user_uuid")
// 	if !exists {
// 		c.JSON(http.StatusForbidden, gin.H{"error": "user_uuid not found"})
// 		return
// 	}

// 	userRole := roleVal.(role.Role)
// 	userUUID := userUUIDVal.(string)

// 	var requests []ds.SoftwareRequest
// 	var err error

// 	if userRole == role.Moderator {
// 		// Модератор видит все заявки
// 		requests, err = h.Repository.GetRequestsFiltered(status, startDate, endDate)
// 	} else {
// 		// Обычный пользователь видит только свои заявки
// 		requests, err = h.Repository.GetRequestsFilteredByUser(userUUID, status, startDate, endDate)
// 	}

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
// 		return
// 	}

// 	// --- Дальше всё остаётся как у тебя ---
// 	type SoftwareResp struct {
// 		SoftwareID  uint   `json:"software_id"`
// 		Title       string `json:"title"`
// 		InstallTime string `json:"install_time"`
// 		FinalTime   string `json:"final_time"`
// 	}

// 	type RequestResp struct {
// 		RequestID      uint           `json:"request_id"`
// 		Status         string         `json:"status"`
// 		CreateDt       string         `json:"create_dt"`
// 		UpdateDt       string         `json:"update_dt"`
// 		FinishDt       string         `json:"finish_dt"`
// 		CreatorID      uint           `json:"creator_id"`
// 		ModeratorID    uint           `json:"moderator_id"`
// 		Phone          string         `json:"phone"`
// 		CompletedCount int64          `json:"completed_count"`
// 		TotalSoftwares int            `json:"total_softwares"`
// 		Softwares      []SoftwareResp `json:"softwares"`
// 	}

// 	var response []RequestResp

// 	for _, r := range requests {
// 		completedCount, _ := h.Repository.CountCompletedSoftware(r.RequestID)

// 		phone := ""
// 		if r.Phone.Valid {
// 			phone = r.Phone.String
// 		}

// 		finishDt := ""
// 		if r.FinishDt.Valid {
// 			finishDt = r.FinishDt.Time.Format(time.RFC3339)
// 		}

// 		installationTimes, err := h.Repository.GetInstallationTimesWithSoftware(r.RequestID)
// 		if err != nil {
// 			logrus.Errorf("Failed to get installation times for request %d: %v", r.RequestID, err)
// 			continue
// 		}

// 		var softwares []SoftwareResp
// 		for _, inst := range installationTimes {
// 			installTime := ""
// 			if inst.InstallTime != nil {
// 				installTime = inst.InstallTime.Format(time.RFC3339)
// 			}

// 			finalTime := ""
// 			if inst.FinalTime != nil {
// 				finalTime = inst.FinalTime.Format(time.RFC3339)
// 			}

// 			softwares = append(softwares, SoftwareResp{
// 				SoftwareID:  inst.SoftwareId,
// 				Title:       inst.Software.Title,
// 				InstallTime: installTime,
// 				FinalTime:   finalTime,
// 			})
// 		}

// 		response = append(response, RequestResp{
// 			RequestID:      r.RequestID,
// 			Status:         r.Status,
// 			CreateDt:       r.CreateDt.Format(time.RFC3339),
// 			UpdateDt:       r.UpdateDt.Format(time.RFC3339),
// 			FinishDt:       finishDt,
// 			CreatorID:      r.CreatorID,
// 			ModeratorID:    r.ModeratorID,
// 			Phone:          phone,
// 			CompletedCount: completedCount,
// 			TotalSoftwares: len(softwares),
// 			Softwares:      softwares,
// 		})
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"total_requests": len(response),
// 		"requests":       response,
// 	})
// }

// GetRequestByID godoc
// @Summary Get request by ID
// @Description Get details of a specific request
// @Tags Requests
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/request/{id} [get]
func (h *Handler) GetRequestByID(c *gin.Context) {
	// 1️⃣ UUID текущего пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем userID
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}
	userID := currentUser.UserID

	// 3️⃣ Получаем ID заявки из URL
	idStr := c.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// 4️⃣ Получаем заявку
	request, err := h.Repository.GetRequestWithSoftware(uint(requestID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	// 5️⃣ Проверяем, что заявка принадлежит текущему пользователю
	if request.CreatorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not have access to this request"})
		return
	}

	// 6️⃣ Приводим phone и finish_dt к читаемому виду
	phone := ""
	if request.Phone.Valid {
		phone = request.Phone.String
	}

	finishDt := ""
	if request.FinishDt.Valid {
		finishDt = request.FinishDt.Time.Format(time.RFC3339)
	}

	// 7️⃣ Получаем программы в заявке с InstallationTime
	installationTimes, err := h.Repository.GetInstallationTimesWithSoftware(request.RequestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get softwares"})
		return
	}

	type SoftwareResp struct {
		SoftwareID  uint   `json:"software_id"`
		Title       string `json:"title"`
		InstallTime string `json:"install_time"`
		FinalTime   string `json:"final_time"`
	}

	var softwares []SoftwareResp
	for _, inst := range installationTimes {
		installTime := ""
		if inst.InstallTime != nil {
			installTime = inst.InstallTime.Format(time.RFC3339)
		}

		finalTime := ""
		if inst.FinalTime != nil {
			finalTime = inst.FinalTime.Format(time.RFC3339)
		}

		softwares = append(softwares, SoftwareResp{
			SoftwareID:  inst.SoftwareId,
			Title:       inst.Software.Title,
			InstallTime: installTime,
			FinalTime:   finalTime,
		})
	}

	// 8️⃣ Формируем итоговый ответ
	response := gin.H{
		"request_id":     request.RequestID,
		"status":         request.Status,
		"create_dt":      request.CreateDt.Format(time.RFC3339),
		"update_dt":      request.UpdateDt.Format(time.RFC3339),
		"finish_dt":      finishDt,
		"creator_id":     request.CreatorID,
		"moderator_id":   request.ModeratorID,
		"phone":          phone,
		"software_count": len(softwares),
		// "softwares":      softwares,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRequestFields godoc
// @Summary Update request fields
// @Description Update fields of user's own request
// @Tags Requests
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param fields body map[string]interface{} true "Fields to update"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests/{id}/fields [put]
func (h *Handler) UpdateRequestFields(c *gin.Context) {
	// 1️⃣ Получаем UUID пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем userID из базы
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}
	userID := currentUser.UserID

	// 3️⃣ Получаем ID заявки из пути
	reqID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// 4️⃣ Проверяем, что заявка принадлежит текущему пользователю
	owned, err := h.Repository.IsRequestOwnedByUser(uint(reqID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !owned {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this request"})
		return
	}

	// 5️⃣ Считываем обновления из JSON
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 6️⃣ Обновляем поля заявки
	if err := h.Repository.UpdateRequestFields(uint(reqID), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// FormRequest godoc
// @Summary Form draft request
// @Description Convert user's draft request into active request
// @Tags Requests
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests/form [put]
func (h *Handler) FormRequest(c *gin.Context) {
	// 1️⃣ Получаем UUID пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем userID из базы
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}
	userID := currentUser.UserID

	// 3️⃣ Получаем черновую заявку текущего пользователя
	draftRequest, err := h.Repository.GetDraftRequestByUser(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "draft request not found"})
		return
	}

	// 4️⃣ Формируем заявку
	if err := h.Repository.FormRequestByDraft(draftRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "formed"})
}

// DeleteRequest godoc
// @Summary Delete request
// @Description Soft delete user's own request
// @Tags Requests
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests/{id} [delete]
func (h *Handler) DeleteRequest(c *gin.Context) {
	// 1️⃣ Получаем UUID пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем userID из базы
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}
	userID := currentUser.UserID

	// 3️⃣ Получаем ID заявки из пути
	idStr := c.Param("id")
	reqID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// 4️⃣ Проверяем, что заявка принадлежит текущему пользователю
	owned, err := h.Repository.IsRequestOwnedByUser(uint(reqID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !owned {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this request"})
		return
	}

	// 5️⃣ Выполняем удаление (статус deleted)
	if err := h.Repository.DeleteRequest(uint(reqID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "request deleted"})
}

// UpdateRequestStatusByModerator godoc
// @Summary Update request status
// @Description Moderator updates request status to completed or rejected
// @Tags Requests
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param status path string true "Status: completed/rejected"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests/{id}/status/{status} [put]
func (h *Handler) UpdateRequestStatusByModerator(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 64)
	status := c.Param("status")

	// Разрешённые статусы
	if status != "completed" && status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	moderator := auth.GetFixedUser()

	if err := h.Repository.CompleteOrRejectRequest(uint(id), moderator.UserID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "request " + status,
	})
}

// DeleteInstallationTime godoc
// @Summary Delete installation record
// @Description Delete installation record for software in user's request
// @Tags Requests
// @Accept json
// @Produce json
// @Param requests_id path int true "Request ID"
// @Param software_id path int true "Software ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests/{requests_id}/software/{software_id} [delete]
func (h *Handler) DeleteInstallationTime(c *gin.Context) {
	// 1️⃣ Получаем UUID текущего пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2️⃣ Получаем пользователя (узнаём user_id)
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}

	// 3️⃣ Читаем параметры
	reqStr := c.Param("requests_id")
	softStr := c.Param("software_id")

	requestID, err1 := strconv.ParseUint(reqStr, 10, 64)
	softwareID, err2 := strconv.ParseUint(softStr, 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request or software id"})
		return
	}

	// 4️⃣ Проверка владельца заявки
	owned, err := h.Repository.IsRequestOwnedByUser(uint(requestID), currentUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !owned {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you do not own this request",
		})
		return
	}

	// 5️⃣ Удаляем запись
	if err := h.Repository.DeleteInstallationTime(uint(requestID), uint(softwareID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "installation record deleted"})
}

// UpdateInstallationTime godoc
// @Summary Update installation time
// @Description Update installation time for software in user's request
// @Tags Requests
// @Accept json
// @Produce json
// @Param request_id path int true "Request ID"
// @Param software_id path int true "Software ID"
// @Param install_time body string true "Install time in format 2006-01-02 15:04"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/requests/{request_id}/software/{software_id}/time [put]
func (h *Handler) UpdateInstallationTime(c *gin.Context) {
	// 1. UUID пользователя из JWT
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid"})
		return
	}

	// 2. Получаем пользователя по UUID, чтобы узнать его user_id
	currentUser, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot load current user"})
		return
	}

	// 3. Получаем параметры
	reqStr := c.Param("request_id")
	softStr := c.Param("software_id")

	requestID, err1 := strconv.ParseUint(reqStr, 10, 64)
	softwareID, err2 := strconv.ParseUint(softStr, 10, 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request or software id"})
		return
	}

	// 4. Проверяем, что заявка принадлежит текущему пользователю
	owned, err := h.Repository.IsRequestOwnedByUser(uint(requestID), currentUser.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !owned {
		c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this request"})
		return
	}

	// 5. Парсим JSON
	var payload struct {
		InstallTime string `json:"install_time"` // формат "2006-01-02 15:04"
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	parsedTime, err := time.Parse("2006-01-02 15:04", payload.InstallTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	// 6. Пытаемся обновить в БД
	if err := h.Repository.UpdateInstallationTime(uint(requestID), uint(softwareID), parsedTime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "installation time updated"})
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get profile info of current user
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} dto.UsersDTO
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/user/profile [get]
func (h *Handler) GetUserProfile(c *gin.Context) {
	// 1. Достаём UUID текущего пользователя из контекста
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}

	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid format"})
		return
	}

	// 2. Получаем пользователя из БД по UUID
	u, err := h.Repository.GetUserByUUID(userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3. Отдаём профиль
	c.JSON(http.StatusOK, u)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update login/password of current user
// @Tags Users
// @Accept json
// @Produce json
// @Param user body object{login=string,password=string} true "User info"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/user/profile [put]
func (h *Handler) UpdateUserProfile(c *gin.Context) {
	// 1. Получаем UUID текущего пользователя из контекста
	userUUIDRaw, ok := c.Get("user_uuid")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user uuid missing in context"})
		return
	}
	userUUID, ok := userUUIDRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user uuid format"})
		return
	}

	// 2. Считываем входные данные
	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 3. Обновляем пользователя по UUID
	if err := h.Repository.UpdateUserByUUID(userUUID, input.Login, input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

// func generateHashString(s string) string {
// 	h := sha1.New()
// 	h.Write([]byte(s))
// 	return hex.EncodeToString(h.Sum(nil))
// }

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type registerReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type registerResp struct {
	Ok bool `json:"ok"`
}

// Login godoc
// @Summary User login
// @Description Login user and get JWT token
// @Tags User
// @Accept json
// @Produce json
// @Param credentials body loginReq true "Login info"
// @Success 200 {object} loginResp
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req loginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	user, err := h.Repository.GetUserByLogin(req.Login)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "user not found"})
		return
	}

	if user.Password != req.Password {
		c.JSON(http.StatusForbidden, gin.H{"error": "wrong password"})
		return
	}
	// Генерация JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ds.JWTClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "bitop-admin",
		},
		UserUUID: user.UUID.String(),
		Role:     user.Role,
	})

	strToken, err := token.SignedString([]byte(h.Config.JWT.Token))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create token"})
		return
	}

	c.JSON(http.StatusOK, loginResp{
		ExpiresIn:   3600,
		AccessToken: strToken,
		TokenType:   "Bearer",
	})
}

// RegisterUser godoc
// @Summary Register new user
// @Description Create new user account
// @Tags User
// @Accept json
// @Produce json
// @Param user body registerReq true "User info"
// @Success 200 {object} registerResp
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/user/register [post]
func (h *Handler) RegisterUser(c *gin.Context) {
	var req registerReq

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	if req.Login == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "login and password required"})
		return
	}

	user := ds.Users{
		UUID:     uuid.New(),
		Login:    req.Login,
		Password: req.Password,
		Role:     role.Creater, // по умолчанию
	}

	if err := h.Repository.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, registerResp{Ok: true})
}

// Logout godoc
// @Summary User logout
// @Description Logout current user, add JWT to blacklist
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security BearerAuth
// @Router /api/user/logout [post]
func (h *Handler) Logout(gCtx *gin.Context) {

	const jwtPrefix = "Bearer "

	// получаем заголовок
	jwtStr := gCtx.GetHeader("Authorization")
	if !strings.HasPrefix(jwtStr, jwtPrefix) { // если нет префикса то нас дурят!
		gCtx.AbortWithStatus(http.StatusBadRequest) // отдаем что нет доступа

		return // завершаем обработку
	}

	// отрезаем префикс
	jwtStr = jwtStr[len(jwtPrefix):]

	_, err := jwt.ParseWithClaims(jwtStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.Config.JWT.Token), nil
	})
	if err != nil {
		gCtx.AbortWithError(http.StatusBadRequest, err)
		log.Println(err)

		return
	}

	// сохраняем в блеклист редиса
	err = h.Redis.WriteJWTToBlacklist(gCtx.Request.Context(), jwtStr, h.Config.JWT.ExpiresIn)
	if err != nil {
		gCtx.AbortWithError(http.StatusInternalServerError, err)

		return
	}

	gCtx.Status(http.StatusOK)
}
