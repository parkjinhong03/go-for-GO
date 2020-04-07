package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

type customClaim struct {
	UserId uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWTString(userId uint) string {
	priv, err := ioutil.ReadFile("./jwt_key.priv")
	if err != nil {
		log.Fatal(err)
	}

	claims := customClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(priv)
	if err != nil {
		log.Fatal(err)
	}
	return ss
}

func ParseCustomClaimsFormJWT(tokenString string) (*customClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaim{}, func(token *jwt.Token) (interface{}, error) {
		return ioutil.ReadFile("./jwt_key.priv")
	})
	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(*customClaim); ok && token.Valid {
		return claims, nil
	}

	return &customClaim{}, fmt.Errorf("can't parse custom claims from JWT")
}
