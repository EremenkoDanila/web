package handler

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"

    "github.com/sirupsen/logrus"
    "lab1/internal/s3"
	"lab1/internal/app/usercontext"
    "lab1/internal/app/ds"
)


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

func (h *Handler) AddToDraftRequest(c *gin.Context) {
    softwareIDStr := c.Param("software_id")
    softwareID, err := strconv.ParseUint(softwareIDStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid software ID"})
        return
    }

    user := auth.GetFixedUser()
    userID := user.GetUserID()

    software, err := h.Repository.GetSoftwareByID(uint(softwareID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Software not found"})
        return
    }

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

    err = h.Repository.AddSoftwareToRequest(draftRequest.RequestID, uint(softwareID))
    if err != nil {
        logrus.Errorf("Failed to add software to request: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add software to draft request"})
        return
    }

    installationTimes, err := h.Repository.GetInstallationTimesWithSoftware(draftRequest.RequestID)
    if err != nil {
        logrus.Errorf("Failed to get installation times: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve installation details"})
        return
    }

    var phone string
    if draftRequest.Phone.Valid {
        phone = draftRequest.Phone.String
    }

    finishDt := ""
    if draftRequest.FinishDt.Valid {
        finishDt = draftRequest.FinishDt.Time.Format(time.RFC3339)
    }

    // --- Формируем список программ с InstallTime и FinalTime ---
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













func (h *Handler) GetUserCart(c *gin.Context) {
    userID := auth.GetFixedUser().GetUserID()

    request, err := h.Repository.GetDraftRequestByUser(userID)
    if err != nil {
        c.JSON(http.StatusOK, gin.H{"request_id": nil, "software_count": 0})
        return
    }

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

    // Получаем заявки
    requests, err := h.Repository.GetRequestsFiltered(status, startDate, endDate)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get requests"})
        return
    }

    type SoftwareResp struct {
        SoftwareID  uint   `json:"software_id"`
        Title       string `json:"title"`
        InstallTime string `json:"install_time"`
        FinalTime   string `json:"final_time"`
    }

    type RequestResp struct {
        RequestID       uint            `json:"request_id"`
        Status          string          `json:"status"`
        CreateDt        string          `json:"create_dt"`
        UpdateDt        string          `json:"update_dt"`
        FinishDt        string          `json:"finish_dt"`
        CreatorID       uint            `json:"creator_id"`
        ModeratorID     uint            `json:"moderator_id"`
        Phone           string          `json:"phone"`
        CompletedCount  int64           `json:"completed_count"`
        TotalSoftwares  int             `json:"total_softwares"`
        Softwares       []SoftwareResp  `json:"softwares"`
    }

    var response []RequestResp

    for _, r := range requests {
        // Получаем количество завершённых программ
        completedCount, _ := h.Repository.CountCompletedSoftware(r.RequestID)

        // Приводим phone и finish_dt в читаемый вид
        phone := ""
        if r.Phone.Valid {
            phone = r.Phone.String
        }

        finishDt := ""
        if r.FinishDt.Valid {
            finishDt = r.FinishDt.Time.Format(time.RFC3339)
        }

        // Получаем программы в заявке с InstallationTime
        installationTimes, err := h.Repository.GetInstallationTimesWithSoftware(r.RequestID)
        if err != nil {
            logrus.Errorf("Failed to get installation times for request %d: %v", r.RequestID, err)
            continue
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

        // Добавляем заявку в итоговый массив
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
            TotalSoftwares: len(softwares),
            Softwares:      softwares,
        })
    }

    // Итоговый JSON
    c.JSON(http.StatusOK, gin.H{
        "total_requests": len(response),
        "requests":       response,
    })
}




func (h *Handler) GetRequestByID(c *gin.Context) {
	idStr := c.Param("id")
	requestID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	// Получаем заявку
	request, err := h.Repository.GetRequestWithSoftware(uint(requestID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	// Приводим phone и finish_dt в читаемый вид
	phone := ""
	if request.Phone.Valid {
		phone = request.Phone.String
	}

    finishDt := ""
    if request.FinishDt.Valid {
        finishDt = request.FinishDt.Time.Format(time.RFC3339)
    }

	// Получаем программы в заявке с InstallationTime
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

	// Формируем итоговый ответ
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
		"softwares":      softwares,
	}

	c.JSON(http.StatusOK, response)
}





// PUT /requests/:id/fields
func (h *Handler) UpdateRequestFields(c *gin.Context) {
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.Repository.UpdateRequestFields(uint(id), updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// PUT /requests/:id/form
func (h *Handler) FormRequest(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    if err := h.Repository.FormRequest(uint(id)); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "formed"})
}




func (h *Handler) DeleteRequest(c *gin.Context) {
idStr := c.Param("id")
id, _ := strconv.ParseUint(idStr, 10, 64)

if err := h.Repository.DeleteRequest(uint(id)); err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}

c.JSON(http.StatusOK, gin.H{"message": "request deleted"})

}



// PUT /requests/:id/status/:status — завершить или отклонить заявку модератором
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









func (h *Handler) DeleteInstallationTime(c *gin.Context) {
    reqStr := c.Param("requests_id")
    softStr := c.Param("software_id")

    requestID, err1 := strconv.ParseUint(reqStr, 10, 64)
    softwareID, err2 := strconv.ParseUint(softStr, 10, 64)

    if err1 != nil || err2 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request or software id"})
        return
    }

    if err := h.Repository.DeleteInstallationTime(uint(requestID), uint(softwareID)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "installation record deleted"})
}



func (h *Handler) UpdateInstallationTime(c *gin.Context) {
    reqStr := c.Param("request_id")
    softStr := c.Param("software_id")

    requestID, err1 := strconv.ParseUint(reqStr, 10, 64)
    softwareID, err2 := strconv.ParseUint(softStr, 10, 64)

    if err1 != nil || err2 != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request or software id"})
        return
    }

    var payload struct {
        InstallTime string `json:"install_time"` // формат: "2006-01-02 15:04"
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

    if err := h.Repository.UpdateInstallationTime(uint(requestID), uint(softwareID), parsedTime); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "installation time updated"})

}













func (h *Handler) RegisterUser(c *gin.Context) {
    var payload struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&payload); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
        return
    }

    user := ds.Users{
        Login:    payload.Login,
        Password: payload.Password,
    }

    if err := h.Repository.CreateUser(&user); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user registered", "user_id": user.UserID})
}
















func (h *Handler) GetUserProfile(c *gin.Context) {
    user := auth.GetFixedUser() // заменить на получение из сессии/токена

    u, err := h.Repository.GetUserByID(user.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, u)
}

// PUT /user — обновление данных текущего пользователя (личный кабинет)
func (h *Handler) UpdateUserProfile(c *gin.Context) {
    user := auth.GetFixedUser()

    var input struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.Repository.UpdateUser(user.UserID, input.Login, input.Password); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}



// POST /auth/login — аутентификация пользователя
func (h *Handler) Login(c *gin.Context) {
    var input struct {
        Login    string `json:"login"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.Repository.Authenticate(input.Login, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    // token, err := auth.CreateSession(user.UserID)
    // if err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    //     return
    // }

    c.JSON(http.StatusOK, gin.H{"message": "login successful", "user": user.UserID})
}



// POST /auth/logout — деавторизация пользователя
func (h *Handler) Logout(c *gin.Context) {
    // user := auth.GetFixedUser()
    // if err := auth.DeleteSession(user.UserID); err != nil {
    //     c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    //     return
    // }

    c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}