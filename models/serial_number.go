package models

type SL_numbers struct {
	Id    int `gorm:"primaryKey,not null,unique"`
	Count int `gorm:"not null"`
}

func (e *SL_numbers) TableName() string {
	return "SL_numbers"
}
