package generate

import (
	"crypto/rsa"
	"crypto/x509"
	"fmt"
	"time"
)

func KeyAndCertificate(
	template, parentCert *x509.Certificate,
	parentKey *rsa.PrivateKey,
	keyPath, certificatePath, password string,
	duration time.Duration,
) (*rsa.PrivateKey, *x509.Certificate) {

	priv := PrivateKey()
	err := SavePrivateKey(priv, keyPath, []byte(password))
	if err != nil {
		panic(err)
	}
	certData := X509Certificate(priv, template, duration, parentKey, parentCert)
	err = SaveX509Certificate(certData, certificatePath)
	if err != nil {
		panic(err)
	}

	// x509.ParseCertificate 함수를 이용하여 바이트 배열로 된 인증서 데이터를 x509.Certificate 객체로 디코딩할 수 있다.
	cert, err := x509.ParseCertificate(certData)
	if err != nil {
		panic(fmt.Sprint("Unable to parse certificate: ", err))
	}

	return priv, cert
}