package generate

import (
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"fmt"
)

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