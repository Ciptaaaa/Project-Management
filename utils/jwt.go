package utils

import (
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