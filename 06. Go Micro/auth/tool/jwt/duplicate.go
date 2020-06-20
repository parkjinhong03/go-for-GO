package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"time"
)

func GenerateDuplicateCertJWT(userId string, d time.Duration) (ss string, err error) {
	priv, err := ioutil.ReadFile("./jwt_key.priv")
	if err != nil {
		return
	}

	claims := duplicateCertClaim{
		UserId:         userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(d).Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = t.SignedString(priv)
	return
}