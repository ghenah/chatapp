package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ghenah/chatapp/pkg/chatapp"
	"github.com/ghenah/chatapp/pkg/dsgorm"
	"github.com/ghenah/chatapp/pkg/httpserver"
	"github.com/joho/godotenv"
)

// @title Chatapp
// @version v0.1.0
// @description A chat app.

// @host localhost:8000
// @BasePath /

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// Set up a database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s%s",
		os.Getenv("MYSQL_DB_USER"),
		os.Getenv("MYSQL_DB_PASSWORD"),
		os.Getenv("MYSQL_DB_HOST"),
		os.Getenv("MYSQL_DB_DATABASE"),
		os.Getenv("MYSQL_DB_DATABASE_SETTINGS"),
	)
	db, err := dsgorm.Init(dsn)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	chatApp, err := chatapp.Init(&chatapp.ChatAppConfig{UsersDS: dsgorm.GetDataStore()})
	if err != nil {
		log.Fatal(err)
	}

	// Start the HTTP server
	serverConfig := httpserver.ServerConfig{
		AppAddressHostname:       os.Getenv("APP_ADDRESS_HOSTNAME"),
		AppAddressPort:           os.Getenv("APP_ADDRESS_PORT"),
		JWTSecretKey:             os.Getenv("APP_JWT_SECRET_KEY"),
		JWTRefreshTokenSecretKey: os.Getenv("APP_JWT_REFRESH_TOKEN_SECRET_KEY"),
		JWTWebSocketSecretKey:    os.Getenv("APP_JWT_WEB_SOCKET_SECRET_KEY"),
		AppWsOriginSchema:        os.Getenv("APP_WS_ORIGIN_SCHEMA"),
		AppWsOriginDomain:        os.Getenv("APP_WS_ORIGIN_DOMAIN"),
		AppWsOriginPort:          os.Getenv("APP_WS_ORIGIN_PORT"),
		// Pass the data store handle
		DS:      dsgorm.GetDataStore(),
		ChatApp: chatApp,
	}
	httpserver.StartServer(serverConfig)
}
