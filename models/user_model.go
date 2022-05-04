package models

import (
	"time"
)

type SL_USER struct {
	User_id             string `gorm:"primaryKey,not null,unique"`
	Phone_Number        string `gorm:"not null"`
	Pin                 string
	Pin_activated       bool
	Active_device_id    int `gorm:"not null"`
	User_status         bool
	Privacy_policy      bool
	Consent_time        *time.Time
	Created_at          time.Time
	Updated_at          time.Time
	Pin_updated_at      *time.Time
	Pin_changed_on      *time.Time
	Pin_change_error    string
	Otp_verified_status bool
}

func (e *SL_USER) TableName() string {
	return "SL_USER"
}
