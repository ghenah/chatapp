package dsgorm

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var dataStore *DataStoreGORM

// Init returns the pointer to the database connection (to allow the caller
// deferring the db.Close() call) and an error.
func Init(dsn string) (*gorm.DB, error) {
	var (
		err error
		db  *gorm.DB
	)

	// Connect to the database and add the connection to the data store handle
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel: logger.Silent,
		},
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return db, err
	}

	err = db.AutoMigrate(&User{})
	if err == nil {
		dataStore = &DataStoreGORM{db: db}
	}

	return db, err
}
