package ds

import (
	"github.com/google/uuid"
	"lab1/internal/app/role"
)

type Users struct {
	UserID       uint      `gorm:"primaryKey;autoIncrement" json:"user_id"`
	UUID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"uuid"`
	Login        string    `gorm:"type:varchar(32);unique;not null" json:"login"`
	Password     string    `gorm:"type:varchar(32);not null" json:"password"` // <-- исправлено
	Role         role.Role `gorm:"type:int" json:"role"`
	ModeratorFlg bool      `gorm:"type:boolean;default:false" json:"moderator_flg"`

	CreatedRequests   []SoftwareRequest `gorm:"foreignKey:CreatorID"`
	ModeratedRequests []SoftwareRequest `gorm:"foreignKey:ModeratorID"`
}
