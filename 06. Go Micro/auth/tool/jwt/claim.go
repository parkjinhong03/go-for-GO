package jwt

import "github.com/dgrijalva/jwt-go"

type duplicateCertClaim struct {
	UserId string `json:"user_id"`
	Email string `json:"email"`
	jwt.StandardClaims
}