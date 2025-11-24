package ds

import (
	"github.com/golang-jwt/jwt"
	"lab1/internal/app/role"
	
)

type JWTClaims struct {
	jwt.StandardClaims
	UserUUID string     `json:"user_uuid"`
	Role     role.Role  `json:"role"`
	Scopes   []string   `json:"scopes"`
}
