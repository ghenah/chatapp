package httpserver

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecretKey string
)

func getJWTSecret() string {
	return jwtSecretKey
}

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

// generateUserSession given the ID of a user, returns the access token
// and an error
func generateUserSession(userID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	accessToken, err := generateAccessToken(userID, expirationTime, []byte(getJWTSecret()))
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// generateAccessToken returns an access token associated with the user ID
// and an error.
func generateAccessToken(userID uint, expirationTime time.Time, secret []byte) (string, error) {
	claims := &Claims{
		UserID: userID,
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
