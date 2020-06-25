package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

func GenerateDuplicateCertJWT(userId, email string, d time.Duration) (ss string, err error) {
	priv, err := ioutil.ReadFile("/Users/parkjinhong/Desktop/go-for-GO/06. Go Micro/auth/tool/jwt/jwt_key.priv")
	if err != nil {
		return
	}

	claims := duplicateCertClaim{
		UserId: userId,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(d).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = t.SignedString(priv)
	return
}

func GenerateDuplicateCertJWTNoReturnErr(userId, email string, d time.Duration) (ss string) {
	priv, err := ioutil.ReadFile("/Users/parkjinhong/Desktop/go-for-GO/06. Go Micro/auth/tool/jwt/jwt_key.priv")
	if err != nil {
		log.Fatal(err)
	}

	claims := duplicateCertClaim{
		UserId: userId,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(d).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = t.SignedString(priv)
	if err != nil {
		log.Fatal(err)
	}
	return
}