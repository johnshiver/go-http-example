package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/johnshiver/asapp_challenge/config"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// CreateToken ...
// Creates new JWT for userID. Nothing is persisted to database
func CreateToken(userID int64) (string, error) {

	c := config.Get()
	expirationTime := time.Now().Add(c.JwtExpiration)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			// library automatically checks ExpiresAt against current time for verification
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(c.JwtSecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func VerifyToken(rawToken string) (*Claims, error) {
	c := config.Get()
	tkn, err := jwt.ParseWithClaims(rawToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(c.JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := tkn.Claims.(*Claims); ok && tkn.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}
