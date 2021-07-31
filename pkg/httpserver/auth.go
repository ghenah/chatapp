package httpserver

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecretKey             string
	jwtRefreshTokenSecretKey string
	jwtWebSocketSecretKey    string
)

func getJWTSecret() string {
	return jwtSecretKey
}
func getRefreshTokenSecret() string {
	return jwtRefreshTokenSecretKey
}
func getJWTWebSocketSecret() string {
	return jwtSecretKey
}

type Claims struct {
	UserID   uint
	Username string
	jwt.StandardClaims
}

// generateUserSession given the ID and the username of a user, returns
// an access token and an error
func generateUserSession(userID uint, username string) (string, error) {
	expirationTime := time.Now().Add(3 * time.Minute)
	accessToken, err := generateAccessToken(userID, username, expirationTime, []byte(getJWTSecret()))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// generateRefreshToken given the ID and the username of a user, returns
// an access token and an error
func generateRefreshToken(userID uint, username string) (string, error) {
	expirationTime := time.Now().Add(72 * time.Hour)
	refreshToken, err := generateAccessToken(userID, username, expirationTime, []byte(getRefreshTokenSecret()))
	if err != nil {
		return "", err
	}

	return refreshToken, nil
}

// generateWebSocketTIcket given the ID and the username of a user, returns
// a web socket ticket and an error
func generateWebSocketTicket(userID uint, username string) (string, error) {
	expirationTime := time.Now().Add(20 * time.Second)
	wsTicket, err := generateAccessToken(userID, username, expirationTime, []byte(getJWTWebSocketSecret()))
	if err != nil {
		return "", err
	}

	return wsTicket, nil
}

// generateAccessToken returns an access token associated with the user ID
// and an error.
func generateAccessToken(userID uint, username string, expirationTime time.Time, secret []byte) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
