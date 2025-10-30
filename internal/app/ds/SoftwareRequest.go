package ds

import (
	"database/sql"
	"time"
)

type SoftwareRequest struct {
	RequestID   uint           `gorm:"primaryKey;autoIncrement" json:"request_id"`
	Status      string         `gorm:"type:varchar(32);not null" json:"status"`
	CreateDt    time.Time      `gorm:"not null" json:"create_dt"`
	UpdateDt    time.Time      `gorm:"not null" json:"update_dt"`
	FinishDt    sql.NullTime   `gorm:"default:null" json:"finish_dt"`
	CreatorID   uint           `gorm:"not null" json:"creator_id"`
	ModeratorID uint           `gorm:"default:null" json:"moderator_id"`
	Phone       sql.NullString `gorm:"type:varchar(32);default:null" json:"phone"`

	Creator   Users `gorm:"foreignKey:CreatorID;references:UserID"`
	Moderator Users `gorm:"foreignKey:ModeratorID;references:UserID"`

	// Для миграции временно закомментировано
	// Installations []InstallationTime `gorm:"foreignKey:RequestID"`
}
