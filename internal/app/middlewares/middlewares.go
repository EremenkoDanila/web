package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	redisv8 "github.com/go-redis/redis/v8" // официальный Redis
	"lab1/internal/app/ds"
	"lab1/internal/app/role"
	myredis "lab1/internal/app/redis" // твой внутренний пакет с клиентом Redis
)

const jwtPrefix = "Bearer "

// WithAuthCheck выполняет JWT-проверку + проверку на blacklist в Redis
func WithAuthCheck(jwtKey string, redisClient *myredis.Client, allowedRoles ...role.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, jwtPrefix) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Authorization header missing or invalid"})
			return
		}

		tokenStr := authHeader[len(jwtPrefix):]

		// ✅ Проверяем, не в blacklist ли токен
		err := redisClient.CheckJWTInBlacklist(context.Background(), tokenStr)
		if err == nil {
			// токен найден в blacklist → запрещаем
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Token is blacklisted"})
			return
		}
		if !errors.Is(err, redisv8.Nil) {
			// внутренняя ошибка Redis
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		// ✅ Проверяем JWT подпись
		token, err := jwt.ParseWithClaims(tokenStr, &ds.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(*ds.JWTClaims)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid token claims"})
			return
		}

		// ✅ Проверка ролей
		if len(allowedRoles) > 0 {
			allowed := false
			for _, r := range allowedRoles {
				if claims.Role == r {
					allowed = true
					break
				}
			}
			if !allowed {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You do not have permission to access this resource"})
				return
			}
		}

		// ✅ Сохраняем данные в контекст
		ctx.Set("user_uuid", claims.UserUUID)
		ctx.Set("role", claims.Role)

		ctx.Next()
	}
}
