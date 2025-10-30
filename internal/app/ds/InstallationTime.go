package ds

import "time"

type InstallationTime struct {
    SoftwareId  uint       `gorm:"primaryKey" json:"software_id"`
    RequestID   uint       `gorm:"primaryKey" json:"request_id"`
    InstallTime *time.Time `gorm:"type:timestamp;default:null" json:"install_time"`
    FinalTime   *time.Time `gorm:"type:timestamp;default:null" json:"final_time"`

    Software Software        `gorm:"foreignKey:SoftwareId;references:SoftwareId"`
    Request  SoftwareRequest `gorm:"foreignKey:RequestID;references:RequestID"`
}
