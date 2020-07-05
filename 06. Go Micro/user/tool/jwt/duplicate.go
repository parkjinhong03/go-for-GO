package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
	"time"
)

func GenerateDuplicateCertJWT(userId, email string, d time.Duration) (ss string, err error) {
	priv, err := ioutil.ReadFile("/Users/parkjinhong/Desktop/go-for-GO/06. Go Micro/auth/tool/jwt/jwt_key.priv")
	if err != nil { return }

	claims := DuplicateCertClaim{
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

	claims := DuplicateCertClaim{
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

func ParseDuplicateCertClaimFromJWT(ss string) (*DuplicateCertClaim, error) {
	token, err := jwt.ParseWithClaims(ss, &DuplicateCertClaim{}, func(t *jwt.Token) (i interface{}, e error) {
		return ioutil.ReadFile("/Users/parkjinhong/Desktop/go-for-GO/06. Go Micro/auth/tool/jwt/jwt_key.priv")
	})
	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*DuplicateCertClaim); ok && token.Valid {
		return claim, nil
	}

	return nil, errors.New("unable to parse duplicate certificate claim from JWT")
}