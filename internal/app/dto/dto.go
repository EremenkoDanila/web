package dto

import "time"

// Для Software
type SoftwareDTO struct {
	SoftwareId      uint    `json:"software_id"`
	DeletedFlg      bool    `json:"deleted_flg"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	FullDescription string  `json:"full_description"`
	ImgURL          string  `json:"img_url"`
	OS              string  `json:"os"`
	Size            float64 `json:"size"`
	Version         string  `json:"version"`
	Functions       string  `json:"functions"`
}

// Для SoftwareRequest
type SoftwareRequestDTO struct {
	RequestID   uint        `json:"request_id"`
	Status      string      `json:"status"`
	CreateDt    time.Time   `json:"create_dt"`
	UpdateDt    time.Time   `json:"update_dt"`
	FinishDt    *time.Time  `json:"finish_dt,omitempty"`
	CreatorID   uint        `json:"creator_id"`
	ModeratorID *uint       `json:"moderator_id,omitempty"`
	Phone       *string     `json:"phone,omitempty"`
}

// Для Users
type UsersDTO struct {
	UserID            uint                  `json:"user_id"`
	UUID              string                `json:"uuid"`
	Login             string                `json:"login"`
	Role              int                   `json:"role"`
	ModeratorFlg      bool                  `json:"moderator_flg"`
	CreatedRequests   []SoftwareRequestDTO  `json:"created_requests,omitempty"`
	ModeratedRequests []SoftwareRequestDTO  `json:"moderated_requests,omitempty"`
}

// Для InstallationTime
type InstallationTimeDTO struct {
	SoftwareId  uint            `json:"software_id"`
	RequestID   uint            `json:"request_id"`
	InstallTime *time.Time      `json:"install_time,omitempty"`
	FinalTime   *time.Time      `json:"final_time,omitempty"`
	Software    SoftwareDTO     `json:"software"`
	Request     SoftwareRequestDTO `json:"request"`
}
