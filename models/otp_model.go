package models

import (
	"time"
)

type SL_OTP_Attempts struct {
	Id               uint   `gorm:"primaryKey,not null,unique"`
	Otp              int    `gorm:"not null"`
	Active_otp       int    `gorm:"not null"`
	User_id          string `gorm:"foreignkey:User_id"`
	Otp_entered_time *time.Time
	Status           string `gorm:"not null"`
	Error_code       string
	Phone_number     string
	Otp_generated    time.Time
	Otp_expiry       time.Time
}

func (e *SL_OTP_Attempts) TableName() string {
	return "SL_OTP_Attempts"
}
