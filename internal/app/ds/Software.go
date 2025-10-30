package ds

type Software struct {
	SoftwareId      uint    `gorm:"primaryKey;autoIncrement" json:"software_id"`
	DeletedFlg      bool    `gorm:"type:boolean;not null;default:false" json:"deleted_flg"`
	Title           string  `gorm:"type:varchar(128);not null" json:"title"`
	Description     string  `gorm:"type:varchar(128)" json:"description"`          
	FullDescription string  `gorm:"type:varchar(128)" json:"full_description"`     
	ImgURL          string  `gorm:"type:varchar(128)" json:"img_url"`              
	OS              string  `gorm:"type:varchar(128);not null" json:"os"`   
	Size            float64 `gorm:"type:varchar(128);not null" json:"size"`
	Version         string  `gorm:"type:varchar(32);not null" json:"version"` 
	Functions       string  `gorm:"type:text" json:"functions"`            

	// Для миграции временно закомментировано
	// Installations []InstallationTime `gorm:"foreignKey:SoftwareId"`
}
