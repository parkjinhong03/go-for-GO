package generate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
	"math/big"
	"time"
)

// 해당 함수는 앞에서 생성한 개인 키와 템플릿을 이용하여 X.509 인증서를 생성하고 반환하는 함수이다.
func X509Certificate(
	key *rsa.PrivateKey,
	template *x509.Certificate,
	duration time.Duration,
	parentKey *rsa.PrivateKey,
	parentCert *x509.Certificate) []byte {

	notBefore := time.Now()
	notAfter := notBefore.Add(duration)

	// NotBefore 필드는 해당 인증서가 특정 시간 전에는 유효하지 않다는 것을 나타낸다.
	template.NotBefore = notBefore
	// NotAfter 필드는 해당 인증서가 특정 시간 후에는 유효하지 않다는 것을 나타낸다.
	template.NotAfter = notAfter

	// 아래 두 구문을 통해 128비트인 임의의 큰 정수를 생성할 수 있다.
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		panic(fmt.Errorf("failed to generate serial number: %v", err))
	}

	// SerialNumber 필드(일련 번호)는 인증서 체인에 대해 고유한 값을 가져야 하는 필드이다.
	template.SerialNumber = serialNumber

	// getSubjectKey 함수를 이용하여 개인 키로부터 subject key를 얻을 수 있다.
	subjectKey, err := getSubjectKey(key)
	if err != nil {
		panic(fmt.Errorf("unable to get subject key: %s", err))
	}

	// 트러스트 체인이 올바르게 작동하는데 필요한 SubjectKey 필드를 작성한다.
	// 해당 인증서가 부모 인증서에 의해 서명된 경우 인증 키 식별자는 부모 인증서의 Subject Key 식별자와 일치한다.
	template.SubjectKeyId = subjectKey

	// 루트 인증서와 같은 경우 자체 서명을 하기 위해서 부모 인증서와 개인 키를 자신의 것으로 대체한다.
	if parentKey == nil {
		parentKey = key
	}
	if parentCert == nil {
		parentCert = template
	}

	// 마지막으로 x509.CreateCertificate 함수를 이용하여 템플릿을 기반으로 실제 인증서를 만들어 반환한다.
	cert, err := x509.CreateCertificate(rand.Reader, template, parentCert, &key.PublicKey, parentKey)
	if err != nil {
		panic(err)
	}
	return cert
}

// 개인 키의 공개 버전의 데이터에 접근하기 위해 선언한 구조체이다.
type subjectPublicKeyInfo struct {
	Algorithm        pkix.AlgorithmIdentifier
	SubjectPublicKey asn1.BitString
}

// 매개변수로 받은 키의 공개 버전을 DER 형식으로 직렬화한 다음, 바이트 부분만 추출하여 반환하는 함수이다.
func getSubjectKey(key *rsa.PrivateKey) ([]byte, error) {
	// x509.MarshalPKIXPublicKey 함수를 이용하여 publicKey를 DER 형식으로 직렬화할 수 있다.
	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal public key: %s", err)
	}

	var subPKI subjectPublicKeyInfo
	// asn1.Unmarshal 함수를 이용하여 위에서 얻은 ASN.1 데이터 구조의 바이트 배열을 subPIK 구조체에 언마샬링 할 수 있다.
	_, err = asn1.Unmarshal(publicKey, &subPKI)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal public key: %s", err)
	}

	// 마지막으로 공개 키 부분의 바이트만 추출하여 반환하면 함수가 완료된다.
	h := sha1.New()
	h.Write(subPKI.SubjectPublicKey.Bytes)
	return h.Sum(nil), nil
}