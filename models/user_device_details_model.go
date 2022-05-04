package models

type SL_User_Device_Details struct {
	Id         int    `gorm:"primaryKey,not null,unique,foreignkey:User_id"`
	User_id    string `gorm:"not null"`
	Location   string `gorm:"not null"`
	Ip_address string `gorm:"not null"`
	Country    string `gorm:"not null"`
	Udid       string `gorm:"not null"`
}

func (e *SL_User_Device_Details) TableName() string {
	return "SL_User_Device_Details"
}
