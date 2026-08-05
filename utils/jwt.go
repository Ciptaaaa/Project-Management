package utils

import (
	"errors"
	"time"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

//generate token jwt

func GenerateToken (userID int64, role,email string, publicID uuid.UUID) (string, error){
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTExpire)
claims :=  jwt.MapClaims{
	"user_id": userID,
	"role":role,
	"public_id":publicID,
	"email":email,
	"exp":time.Now().Add(duration).Unix(),
}

token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

return token.SignedString([]byte(secret))
}
//generate refresh token
func GenerateRefreshToken (userID int64) (string, error){
	secret := config.AppConfig.JWTSecret
	duration, _ := time.ParseDuration(config.AppConfig.JWTRefreshToken)
claims :=  jwt.MapClaims{
	"user_id": userID,
	"exp":time.Now().Add(duration).Unix(),
}

token:= jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

return token.SignedString([]byte(secret))
}

func ParseRefreshToken(tokenStr string) (int64, error) {
	secret := config.AppConfig.JWTSecret

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid refresh token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64) 
	if !ok {
		return 0, errors.New("user_id claim missing")
	}

	return int64(userIDFloat), nil
}