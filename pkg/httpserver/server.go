package httpserver

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/ghenah/chatapp/docs"
	"github.com/ghenah/chatapp/pkg/chatapp"
	"github.com/ghenah/chatapp/pkg/idatastore"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	echoSwagger "github.com/swaggo/echo-swagger"
)

type ServerConfig struct {
	AppAddressHostname       string
	AppAddressPort           string
	AppWsOriginSchema        string
	AppWsOriginDomain        string
	AppWsOriginPort          string
	JWTSecretKey             string
	JWTRefreshTokenSecretKey string
	JWTWebSocketSecretKey    string
	DS                       idatastore.IDataStore
	ChatApp                  *chatapp.ChatApp
}

var e *echo.Echo
var ds idatastore.IDataStore
var ca *chatapp.ChatApp

func StartServer(cfg ServerConfig) {
	ds = cfg.DS
	ca = cfg.ChatApp

	if cfg.JWTSecretKey == "" {
		log.Fatal("no value for JWTSecretKey found in the environment")
	} else if cfg.JWTWebSocketSecretKey == "" {
		log.Fatal("no value for JWTWebSocketSecretKey found in the environment")
	} else if cfg.JWTRefreshTokenSecretKey == "" {
		log.Fatal("no value for JWTRefreshTokenSecretKey found in the environment")
	}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == fmt.Sprintf("%s%s:%s", cfg.AppWsOriginSchema, cfg.AppWsOriginDomain, cfg.AppWsOriginPort)
		},
	}

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
	refreshToken := e.Group("/refresh-token")
	refreshToken.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &Claims{},
		SigningKey:  []byte(getRefreshTokenSecret()),
		TokenLookup: "cookie:refreshToken",
		ContextKey:  "refreshToken",
	}))
	refreshToken.GET("", refreshAccessToken)

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

	// users
	protected.GET("/users/search", userSearch)
	protected.POST("/users/friends", userFriendAdd)
	protected.DELETE("/users/friends", userFriendRemove)
	protected.POST("/users/ignored", userIgnoredAdd)
	protected.DELETE("/users/ignored", userIgnoredRemove)
	protected.GET("/users/profile", getAuthenticatedUserInfo)
	protected.PUT("/users/update/password", userUpdatePassword)
	protected.PUT("/users/update/username", userUpdateUsername)

	// chat
	protected.GET("/chat/ticket", chatGetWSTicket)
	protected.GET("/chat/rooms/search", chatRoomSearch)

	wsProtected := e.Group("/ws")
	wsProtected.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:      &Claims{},
		SigningKey:  []byte(getJWTWebSocketSecret()),
		TokenLookup: "query:ticket",
		ContextKey:  "ticket",
	}))

	wsProtected.GET("/connect", chatConnectionInit)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
}
