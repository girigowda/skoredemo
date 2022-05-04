package models

import (
	"time"
)

type SL_KYC_facematch struct {
	Id         uint    `gorm:"primaryKey,not null,unique"`
	User_id    string  `gorm:"foreignkey:User_id"`
	Status     string  `gorm:"not null"`
	Score      float64 `gorm:"not null"`
	Match      bool    `gorm:"not null"`
	Created_at time.Time
	SL_USER    []SL_USER `gorm:"foreignkey:User_id"`
}

func (e *SL_KYC_facematch) TableName() string {
	return "SL_KYC_facematch"
}
