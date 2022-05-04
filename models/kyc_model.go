package models

import (
	"time"
)

type SL_KYC_Details struct {
	Id                 int    `gorm:"primaryKey,not null,unique"`
	User_id            string `gorm:"foreignkey:User_id"`
	Gender             string `gorm:"not null"`
	Nik_number         string `gorm:"not null"`
	Email              string
	Dob                string
	Profile_image      string
	Created_at         time.Time
	Updated_at         time.Time
	Card_image         string
	Card_upload_staus  bool
	Address            string `gorm:"not null"`
	Village            string
	District           string `gorm:"not null"`
	City               string `gorm:"not null"`
	Province           string
	Family_card_no     int
	Mother_maiden_name string
	Full_name          string `gorm:"not null"`
	First_name         string
	Last_name          string
	Blood_group        string `gorm:"not null"`
	Relegion           string
	Citizenship        string
	Profession         string
	Place_of_birth     string `gorm:"not null"`
	Rt_rw              string `gorm:"not null"`
	Marital_status     string
	Transaction_id     int
	Kyc_approved       bool
	Occupation         string
	ExpiryDate         string `gorm:"not null"`
	Nationality        string
	SL_USER            []SL_USER `gorm:"foreignkey:User_id"`
}

func (e *SL_KYC_Details) TableName() string {
	return "SL_KYC_Details"
}
