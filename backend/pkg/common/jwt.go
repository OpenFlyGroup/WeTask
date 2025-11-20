package common

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// ? InitJWT initializes JWT secret
func InitJWT(secret string) {
	jwtSecret = []byte(secret)
}

// ? Claims represents JWT claims
type Claims struct {
	UserID uint `json:"sub"`
	jwt.RegisteredClaims
}

// ? GenerateToken generates a JWT token
func GenerateToken(userID uint, expirationTime time.Duration, issuedAt time.Time) (string, error) {
	expirationTimePoint := issuedAt.Add(expirationTime)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTimePoint),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ? ValidateToken validates a JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
