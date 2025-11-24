package ds

type Users struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement" json:"user_id"`
	Login        string `gorm:"type:varchar(32);unique;not null" json:"login"`
	Password     string `gorm:"type:varchar(32);not null" json:"password"`
	ModeratorFlg bool   `gorm:"type:boolean;default:false" json:"moderator_flg"`

	CreatedRequests   []SoftwareRequest `gorm:"foreignKey:CreatorID"`
	ModeratedRequests []SoftwareRequest `gorm:"foreignKey:ModeratorID"`
}
