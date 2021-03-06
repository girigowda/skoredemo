package database

import (
	"log"
	"sk-integrated-services/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB    *gorm.DB
	err   error
	DBErr error
)

type Database struct {
	*gorm.DB
}

// Connection create database connection
func Connection() error {
	var db = DB
	dsn := config.DbConfiguration()

	// logMode := viper.GetBool("database.log_mode")
	// loglevel := logger.Silent
	// if logMode {
	// 	loglevel = logger.Info
	// }

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		DBErr = err
		log.Println("Db connection error")
		return err
	}

	// err = db.AutoMigrate(migrationModels...)

	// if err != nil {
	// 	return err
	// }
	DB = db

	return nil

}

// GetDB connection
func GetDB() *gorm.DB {
	return DB
}

// GetDBError connection error
func GetDBError() error {
	return DBErr
}
