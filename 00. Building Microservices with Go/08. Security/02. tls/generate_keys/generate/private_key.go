package generate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

// rsa.GenerateKey 함수를 이용하여 RSA 알고리즘을 사용하는 개인 키를 생성하고 반환하는 함수이다.
func PrivateKey() *rsa.PrivateKey {
	// 첫 번째 매개 변수는 난수를 리턴하는 I/O 리더이고, 두 번째 매개 변수는 사용할 키의 비트 크기이다.
	key, _ := rsa.GenerateKey(rand.Reader, 4096)
	return key
}

func SavePrivateKey(key *rsa.PrivateKey, path string, password []byte) error {
	// x509.MarshalPKCS1PrivateKey 함수는 RSA Private Key를 PKCS#1로 변환(바이트로 마샬링)시켜 반환하는 함수이다.
	b := x509.MarshalPKCS1PrivateKey(key)
	var block *pem.Block
	var err error

	if len(password) > 3 {
		// 만약 비밀번호가 3자리 초과라면, X509.EncryptPEMBlock 함수를 이용하여 DER로 인코딩된 데이터를 주어진 암호로 암호화한 PEM 유형의 블록을 생성한다.
		// 3번째 매개 변수에는 암호화화 할 데이터를, 4번째에는 암호화하는데 사용할 비밀번호를, 5번째에는 암호화하는데 사용할 알고리즘을 작성한다.
		block, err = x509.EncryptPEMBlock(rand.Reader, "RSA PRIVATE KEY", b, password, x509.PEMCipherAES256)
		if err != nil {
			return fmt.Errorf("unable to encrpt key: %s", err)
		}
	} else {
		// 만약 비밀번로가 3가지 이하라면, 내용을 암호화하지 않고 초깃값만 설정한 PEM 유형의 블록을 생성한다.
		block = &pem.Block{Type: "RSA PRIVATE KEY", Bytes: b}
	}

	// os.OpenFile 함수를 이용하여 매개 변수로 받은 경로의 파일에 대한 연결을 생성할 수 있다.
	// 2번째 매개 변수를 통해 읽기 전용(os.O_WRONLY), 존재하지 않다면 생성(os.O_CREATE), 파일 오픈 시 해당 파일의 내용 삭제(os.O_TRUNC) 조건을 줄 수 있다.
	keyOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open key.pem for writing: %v", err)
	}

	err = pem.Encode(keyOut, block)
	if err != nil {
		return fmt.Errorf("failed to write key.pem on file: %v", err)
	}
	_ = keyOut.Close()

	return nil
}