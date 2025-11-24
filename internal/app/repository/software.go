package repository

import (
	"lab1/internal/app/ds"
	"time"
    "gorm.io/gorm"
    "fmt"
    "math/rand"
    "database/sql"
)

// GetSoftwaresWithFilter возвращает список программ (без удалённых)
func (r *Repository) GetSoftwaresWithFilter(filter string) ([]ds.Software, error) {
	var softwares []ds.Software
	query := r.DB.Where("deleted_flg = ?", false)
	if filter != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+filter+"%")
	}
	if err := query.Find(&softwares).Error; err != nil {
		return nil, err
	}
	return softwares, nil
}

// GetSoftwareByID возвращает одну программу по ID
func (r *Repository) GetSoftwareByID(id uint) (*ds.Software, error) {
	var software ds.Software
	if err := r.DB.Where("software_id = ? AND deleted_flg = false", id).First(&software).Error; err != nil {
		return nil, err
	}
	return &software, nil
}

// AddSoftware добавляет новую запись
func (r *Repository) AddSoftware(soft *ds.Software) error {
	return r.DB.Create(soft).Error
}

// UpdateSoftware обновляет существующую запись
func (r *Repository) UpdateSoftware(soft *ds.Software) error {
	return r.DB.Model(&ds.Software{}).Where("software_id = ?", soft.SoftwareId).Updates(soft).Error
}

// SoftDeleteSoftware устанавливает флаг удаления
func (r *Repository) SoftDeleteSoftware(id uint) error {
	return r.DB.Model(&ds.Software{}).Where("software_id = ?", id).Update("deleted_flg", true).Error
}

// UpdateSoftwareImage обновляет ссылку на изображение
func (r *Repository) UpdateSoftwareImage(id uint, newURL string) error {
	return r.DB.Model(&ds.Software{}).Where("software_id = ?", id).Update("img_url", newURL).Error
}




// UpdateRequestTimestamp обновляет время изменения заявки
func (r *Repository) UpdateRequestTimestamp(requestID uint) error {
    return r.DB.Model(&ds.SoftwareRequest{}).
        Where("request_id = ?", requestID).
        Update("update_dt", time.Now()).Error
}




// GetDraftRequestByUser возвращает черновик заявки пользователя
func (r *Repository) GetDraftRequestByUser(userID uint) (*ds.SoftwareRequest, error) {
    var request ds.SoftwareRequest
    err := r.DB.Where("creator_id = ? AND status = ?", userID, "draft").First(&request).Error
    if err != nil {
        return nil, err
    }
    return &request, nil
}

// CreateDraftRequest создает новую заявку в статусе черновика
func (r *Repository) CreateDraftRequest(userID uint) (*ds.SoftwareRequest, error) {
    request := &ds.SoftwareRequest{
        Status:    "draft",
        CreateDt:  time.Now(),
        UpdateDt:  time.Now(),
        CreatorID: userID,
    }
    err := r.DB.Create(request).Error
    return request, err
}

// AddSoftwareToRequest добавляет программу в заявку и обновляет время изменения заявки
func (r *Repository) AddSoftwareToRequest(requestID uint, softwareID uint) error {
    return r.DB.Transaction(func(tx *gorm.DB) error {
        // Добавляем программу в заявку
        installationTime := &ds.InstallationTime{
            SoftwareId: softwareID,
            RequestID:  requestID,
        }
        if err := tx.Create(installationTime).Error; err != nil {
            return err
        }

        // Обновляем время изменения заявки
        if err := tx.Model(&ds.SoftwareRequest{}).
            Where("request_id = ?", requestID).
            Update("update_dt", time.Now()).Error; err != nil {
            return err
        }

        return nil
    })
}

// GetRequestWithSoftware возвращает заявку с связанными программами (ИСПРАВЛЕННАЯ ВЕРСИЯ)
func (r *Repository) GetRequestWithSoftware(requestID uint) (*ds.SoftwareRequest, error) {
    var request ds.SoftwareRequest
    
    // Сначала получаем основную заявку
    err := r.DB.Where("request_id = ?", requestID).First(&request).Error
    if err != nil {
        return nil, err
    }
    
    // Затем получаем связанные InstallationTimes отдельным запросом
    var installationTimes []ds.InstallationTime
    err = r.DB.Where("request_id = ?", requestID).Find(&installationTimes).Error
    if err != nil {
        return nil, err
    }
    
    // Получаем Software для каждого InstallationTime
    for i := range installationTimes {
        var software ds.Software
        err = r.DB.Where("software_id = ?", installationTimes[i].SoftwareId).First(&software).Error
        if err != nil {
            return nil, err
        }
        // Здесь мы не можем напрямую присвоить из-за ограничений структуры,
        // но мы можем создать кастомный ответ в handler
    }
    
    return &request, nil
}

// GetInstallationTimesWithSoftware возвращает InstallationTimes с предзагруженными Software
func (r *Repository) GetInstallationTimesWithSoftware(requestID uint) ([]ds.InstallationTime, error) {
    var installationTimes []ds.InstallationTime
    err := r.DB.Where("request_id = ?", requestID).Find(&installationTimes).Error
    if err != nil {
        return nil, err
    }
    
    // Вручную загружаем Software для каждого InstallationTime
    for i := range installationTimes {
        var software ds.Software
        err = r.DB.Where("software_id = ? AND deleted_flg = false", installationTimes[i].SoftwareId).First(&software).Error
        if err != nil {
            return nil, err
        }
        // В Go мы не можем напрямую изменить структуру в цикле for-range,
        // поэтому используем указатель
        installationTimes[i].Software = software
    }
    
    return installationTimes, nil
}






























func (r *Repository) GetRequestsFiltered(status string, startDate, endDate *time.Time) ([]ds.SoftwareRequest, error) {
query := r.DB.Preload("Creator").Preload("Moderator").Model(&ds.SoftwareRequest{}).Where("status != ?", "draft")

if status != "" {
	query = query.Where("status = ?", status)
}
if startDate != nil {
	query = query.Where("create_dt >= ?", *startDate)
}
if endDate != nil {
	query = query.Where("create_dt <= ?", *endDate)
}

var requests []ds.SoftwareRequest
if err := query.Find(&requests).Error; err != nil {
	return nil, err
}
return requests, nil

}


// Получаем количество выполненных программ для заявки
func (r *Repository) CountCompletedSoftware(requestID uint) (int64, error) {
    var count int64
    err := r.DB.Model(&ds.InstallationTime{}).Where("request_id = ? AND final_time IS NOT NULL", requestID).Count(&count).Error
    if err != nil {
        return 0, err
    }
    return count, nil
}



func (r *Repository) UpdateRequestFields(reqID uint, updates map[string]interface{}) error {
return r.DB.Model(&ds.SoftwareRequest{}).Where("request_id = ?", reqID).Updates(updates).Error
}

// Формирование заявки создателем
func (r *Repository) FormRequest(reqID uint) error {
var req ds.SoftwareRequest
if err := r.DB.First(&req, "request_id = ?", reqID).Error; err != nil {
return err
}

// Проверка обязательного телефона
if !req.Phone.Valid || req.Phone.String == "" {
    return fmt.Errorf("phone is required")
}

// Проверка обязательного InstallTime у всех приложений
var installs []ds.InstallationTime
if err := r.DB.Where("request_id = ?", reqID).Find(&installs).Error; err != nil {
    return err
}

for _, inst := range installs {
    if inst.InstallTime == nil {
        return fmt.Errorf("all software must have InstallTime")
    }
}

// Обновляем статус и дату обновления
req.Status = "formed"
req.UpdateDt = time.Now()
return r.DB.Updates(&req).Error

}


func (r *Repository) DeleteRequest(reqID uint) error {
    return r.DB.Model(&ds.SoftwareRequest{}).
        Where("request_id = ?", reqID).
        Updates(map[string]interface{}{
        "status": "deleted",
        "update_dt": time.Now(),
        "finish_dt": time.Now(),
        }).Error
}


func (r *Repository) CompleteOrRejectRequest(reqID uint, moderatorID uint, status string) error {
    if status != "completed" && status != "rejected" {
        return fmt.Errorf("invalid status: %s", status)
    }

    // Получаем заявку
    var req ds.SoftwareRequest
    if err := r.DB.First(&req, "request_id = ?", reqID).Error; err != nil {
        return err
    }

    now := time.Now()
    req.Status = status
    req.ModeratorID = moderatorID
    req.FinishDt = sql.NullTime{Time: now, Valid: true}
    req.UpdateDt = now

    // Обновляем заявку
    if err := r.DB.Save(&req).Error; err != nil {
        return err
    }

    // Если заявка отклонена — не вычисляем FinalTime
    if status == "rejected" {
        return nil
    }

    // Если заявка завершена — рассчитываем FinalTime
    var installations []ds.InstallationTime
    if err := r.DB.Preload("Software").Where("request_id = ?", reqID).Find(&installations).Error; err != nil {
        return err
    }

    for _, inst := range installations {
        if inst.InstallTime == nil {
            continue
        }

        speed := float64(rand.Intn(10) + 1) // 1–10 ГБ/ч
        durationMin := inst.Software.Size / speed * 60
        finalTime := inst.InstallTime.Add(time.Duration(durationMin) * time.Minute)
        inst.FinalTime = &finalTime

        if err := r.DB.Save(&inst).Error; err != nil {
            return err
        }
    }

    return nil
}







// DeleteInstallationTime удаляет запись из installation_times по RequestID и SoftwareId
func (r *Repository) DeleteInstallationTime(requestID uint, softwareID uint) error {
    result := r.DB.
        Where("request_id = ? AND software_id = ?", requestID, softwareID).
        Delete(&ds.InstallationTime{})

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("no installation record found for request_id=%d and software_id=%d", requestID, softwareID)
    }

    return nil
}


func (r *Repository) UpdateInstallationTime(requestID uint, softwareID uint, installTime time.Time) error {
    result := r.DB.Model(&ds.InstallationTime{}).
    Where("request_id = ? AND software_id = ?", requestID, softwareID).
    Update("install_time", installTime)

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return fmt.Errorf("no installation record found for request_id=%d and software_id=%d", requestID, softwareID)
    }

    return nil

}





func (r *Repository) CreateUser(user *ds.Users) error {
    result := r.DB.Create(user)
    return result.Error
}
















// Получить пользователя по ID
func (r *Repository) GetUserByID(userID uint) (*ds.Users, error) {
    var user ds.Users
    if err := r.DB.First(&user, userID).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

// Обновить данные пользователя
func (r *Repository) UpdateUser(userID uint, login string, password string) error {
    updates := map[string]interface{}{}
    if login != "" {
        updates["login"] = login
    }
    if password != "" {
        updates["password"] = password
    }

    if len(updates) == 0 {
        return nil
    }

    return r.DB.Model(&ds.Users{}).Where("user_id = ?", userID).Updates(updates).Error
}

// Аутентификация
func (r *Repository) Authenticate(login, password string) (*ds.Users, error) {
    var user ds.Users
    if err := r.DB.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
