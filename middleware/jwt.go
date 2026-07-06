package middleware

import (
	"strings"

	"github.com/Ciptaaaa/Project-Management.git/config"
	"github.com/Ciptaaaa/Project-Management.git/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		// Ambil header Authorization
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(ctx, "Missing token", "Authorization header required")
		}

		// Format harus "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return utils.Unauthorized(ctx, "Invalid token format", "Use: Bearer <token>")
		}

		tokenStr := parts[1]

		// Parse & validasi token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Pastikan algoritma HMAC, bukan algorithm confusion attack
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return utils.Unauthorized(ctx, "Invalid or expired token", err.Error())
		}

		// Simpan claims ke context supaya bisa diakses di handler
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return utils.Unauthorized(ctx, "Invalid token claims", "")
		}
		ctx.Locals("user", claims)

		return ctx.Next()
	}
}