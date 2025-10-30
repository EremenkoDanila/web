package repository

import (
	"fmt"
	"lab1/internal/app/ds"
	"time"
)

// Software Queries
func (r *Repository) GetAllSoftware() ([]ds.Software, error) {
	var sw []ds.Software
	err := r.DB.Where("deleted_flg = false").Find(&sw).Error
	return sw, err
}

func (r *Repository) GetSoftwareByID(id uint) (*ds.Software, error) {
	var sw ds.Software
	err := r.DB.Where("software_id = ? AND deleted_flg = false", id).First(&sw).Error
	if err != nil {
		return nil, err
	}
	return &sw, err
}

func (r *Repository) SearchSoftwareByTitle(title string) ([]ds.Software, error) {
	var sw []ds.Software
	err := r.DB.Where("(title ILIKE ? OR description ILIKE ?) AND deleted_flg = false",
		"%"+title+"%", "%"+title+"%").Find(&sw).Error
	return sw, err
}

func (r *Repository) SoftDeleteSoftware(id uint) error {
	return r.DB.Model(&ds.Software{}).Where("software_id = ?", id).Update("deleted_flg", true).Error
}

// Request
func (r *Repository) GetDraftRequest(userID uint) (*ds.SoftwareRequest, error) {
	var req ds.SoftwareRequest
	err := r.DB.Where("creator_id = ? AND status = 'draft'", userID).First(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, err
}

func (r *Repository) CreateDraftRequest(userID uint) (*ds.SoftwareRequest, error) {
	req := ds.SoftwareRequest{
		Status:    "draft",
		CreateDt:  time.Now(),
		UpdateDt:  time.Now(),
		CreatorID: userID,
	}
	err := r.DB.Create(&req).Error
	return &req, err
}

func (r *Repository) AddSoftwareToRequest(requestID, softwareID uint) error {
	install := ds.InstallationTime{
		RequestID:  requestID,
		SoftwareId: softwareID,
		InstallTime: nil,
		FinalTime:   nil,
	}

	err := r.DB.Create(&install).Error
	if err != nil {
		return err
	}

	// обновляем update_dt заявки
	return r.DB.Model(&ds.SoftwareRequest{}).
		Where("request_id = ?", requestID).
		Update("update_dt", time.Now()).
		Error
}


func (r *Repository) GetCartCount(userID uint) int64 {
	var count int64
	req, err := r.GetDraftRequest(userID)
	if err != nil {
		return 0
	}
	r.DB.Model(&ds.InstallationTime{}).Where("request_id = ?", req.RequestID).Count(&count)
	return count
}

func (r *Repository) UpdateInstallTime(requestID, softwareID uint, installTime, finalTime *time.Time) error {
	err := r.DB.Model(&ds.InstallationTime{}).
		Where("request_id = ? AND software_id = ?", requestID, softwareID).
		Updates(map[string]interface{}{
			"install_time": installTime,
			"final_time":   finalTime,
		}).Error

	if err != nil {
		return err
	}

	// обновляем update_dt заявки
	return r.DB.Model(&ds.SoftwareRequest{}).
		Where("request_id = ?", requestID).
		Update("update_dt", time.Now()).
		Error
}


func (r *Repository) DeleteRequestByCursor(requestID, userID uint) error {
	cursor := fmt.Sprintf(
		"UPDATE software_requests SET status='deleted' WHERE request_id=%d AND creator_id=%d",
		requestID, userID)
	return r.DB.Exec(cursor).Error
}

// ===== Новые методы =====
func (r *Repository) GetRequestByID(requestID, userID uint) (*ds.SoftwareRequest, error) {
	var req ds.SoftwareRequest
	err := r.DB.Preload("Creator").
		Where("request_id = ? AND creator_id = ?", requestID, userID).
		First(&req).Error
	return &req, err
}

func (r *Repository) GetInstallationsByRequestID(requestID uint) ([]ds.InstallationTime, error) {
	var installations []ds.InstallationTime
	err := r.DB.Preload("Software").
		Where("request_id = ?", requestID).
		Find(&installations).Error
	return installations, err
}

func (r *Repository) GetInstallationCount(requestID, softwareID uint) int64 {
	var count int64
	r.DB.Model(&ds.InstallationTime{}).
		Where("request_id = ? AND software_id = ?", requestID, softwareID).
		Count(&count)
	return count
}

func (r *Repository) UpdatePhone(requestID, userID uint, phone string) error {
	err := r.DB.Model(&ds.SoftwareRequest{}).
		Where("request_id = ? AND creator_id = ? AND status = 'draft'", requestID, userID).
		Update("phone", phone).Error

	if err != nil {
		return err
	}

	// обновляем update_dt заявки
	return r.DB.Model(&ds.SoftwareRequest{}).
		Where("request_id = ?", requestID).
		Update("update_dt", time.Now()).
		Error
}
