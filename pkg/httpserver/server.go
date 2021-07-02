package httpserver

import (
	_ "github.com/ghenah/chatapp/docs"
	"github.com/ghenah/chatapp/pkg/idatastore"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type ServerConfig struct {
	AppAddressHostname string
	AppAddressPort     string
	DS                 idatastore.IDataStore
}

var e *echo.Echo
var ds idatastore.IDataStore

func StartServer(cfg ServerConfig) {
	ds = cfg.DS

	e = echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${status}\t${method}\t${uri}\n",
	}))
	e.HideBanner = true

	setUpRouter()

	e.Logger.Fatal(e.Start(cfg.AppAddressHostname + ":" + cfg.AppAddressPort))
}

func setUpRouter() {
	auth := e.Group("/auth")

	auth.POST("/signup", userRegister)
	auth.POST("/signin", userAuthencticate)

	api := e.Group("/api/v1")
	protected := api.Group("")
	protected.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &Claims{},
		SigningKey:  []byte(getJWTSecret()),
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ContextKey:  "user",
	}))

	protected.GET("/users/search", userSearch)
	protected.POST("/users/friends", userFriendAdd)
	protected.DELETE("/users/friends", userFriendRemove)
	protected.POST("/users/ignored", userIgnoredAdd)
	protected.DELETE("/users/ignored", userIgnoredRemove)
	protected.PUT("/users/update/password", userUpdatePassword)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
