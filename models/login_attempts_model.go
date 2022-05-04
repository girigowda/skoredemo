package models

import (
	"time"
)

type SL_Login_Attempts struct {
	Id              int    `gorm:"primaryKey,not null,unique"`
	User_id         string `gorm:"foreignkey:User_id"`
	Login_timestamp *time.Time
	Status          string `gorm:"not null"`
	Error_code      string
	Device_id       int `gorm:"not null"`
}

func (e *SL_Login_Attempts) TableName() string {
	return "SL_Login_Attempts"
}
