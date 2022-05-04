package models

import (
	"time"
)

type SL_PIN_Attempts struct {
	Id            int    `gorm:"primaryKey;not null;unique"`
	User_id       string `gorm:"foreignkey:User_id"`
	Status        string `gorm:"not null"`
	Error_code    string
	Pin_generated time.Time `gorm:"not null"`
}

func (e *SL_PIN_Attempts) TableName() string {
	return "SL_PIN_Attempts"
}
