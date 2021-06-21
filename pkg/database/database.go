package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

type Config struct {
	ServerName string
	User       string
	Password   string
	DB         string
}

var GetConnectionString = func(config Config) string {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&multiStatements=true", config.User, config.Password, config.ServerName, config.DB)

	return connectionString
}

var Connector *gorm.DB

func Connect(connectionString string) error {
	var err error

	Connector, err = gorm.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	Connector.AutoMigrate(&Person{})

	log.Println("DB connection successful")
	return nil
}

type Person struct {
	ID        int `json:"id"`
	Username  string `json:"username"`
	Email       string `json:"email"`
}
