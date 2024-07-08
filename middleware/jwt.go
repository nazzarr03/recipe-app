package middleware

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}
