package dbhelpers

import (
	"gorm.io/gorm"
)

func FilterQuery(conn *gorm.DB, table_name string, condition map[string]interface{}, scanType interface{}) *gorm.DB {
	return conn.Table(table_name).Where(condition).Scan(scanType)
}

func InsertQuery(conn *gorm.DB, table_name string, condition map[string]interface{}) *gorm.DB {
	return conn.Table(table_name).Create(condition)
}

func UpdateQuery(conn *gorm.DB, table_name string, condition map[string]interface{}, columns map[string]interface{}) *gorm.DB {
	return conn.Table(table_name).Where(condition).Updates(columns)
}
