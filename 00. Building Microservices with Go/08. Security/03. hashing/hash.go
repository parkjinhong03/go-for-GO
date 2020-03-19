package hashing

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"math/big"
)

// pepper에 대한 배열을 필드로 가지고 있는 Hashing을 하기 위한 객체이다.
type Hash struct {
	peppers []string
}

type SaltConfig struct {
	with 	bool
	byteLen int
}

type PepperConf struct {
	with bool
}

func New(peppers []string) *Hash {
	return &Hash{peppers: peppers}
}

// 솔트, 페퍼의 사용 여부를 판단하고 솔트, 페퍼의 값을 만들어서 해시한 후, 그 값과 솔트의 값을 반환하는 함수이다.
func (h *Hash) GenerateHash(input string, saltConf SaltConfig, pepperConf PepperConf) (hash, salt string) {
	pepper := ""
	if pepperConf.with == true && len(h.peppers) > 0 {
		pepper = h.getRandomPepper()
	}

	salt = ""
	if saltConf.with == true {
		salt = generateRandomSalt(saltConf.byteLen)
	}

	hash = h.createHash(input, salt, pepper)
	return
}

func (h *Hash) Compare(input, salt, hash string, pepperConf PepperConf) bool {
	if pepperConf.with == false {
		return h.createHash(input, salt, "") == hash
	}

	for _, pepper:= range h.peppers {
		if h.createHash(input, salt, pepper) == hash { return true }
	}

	return false
}

// 매개 변수로 받은 세 문자열을 다 합해서 sha256 알고리즘으로 해시한 후, 문자열로 변환시켜 반환하는 함수이다.
func (h *Hash) createHash(input, salt, pepper string) string {
	stringToHash := input + salt + pepper

	sha := sha256.New()
	sha.Write([]byte(stringToHash))

	hash := sha.Sum(nil)
	return hex.EncodeToString(hash)
}

// Hash 객체의 peppers 배열 필드에 저장된 여러 문자열 중 임의로 하나를 골라서 반환하는 함수이다.
func (h *Hash) getRandomPepper() string {
	// 아래 구문을 통해서 h.peppers 배열의 인덱스 갯수에서 임의의 숫자를 반환받을 수 있다.
	max := big.NewInt(int64(len(h.peppers)))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(err)
	}

	// 해당 숫자에 일치하는 인덱스 반환
	return h.peppers[n.Int64()]
}

// 매개변수로 받은 정수의 길이만큼의 바이트 배열의 값을 랜덤으로 대입해서 문자열로 변환하여 반환하는 함수이다.
func generateRandomSalt(byteLen int) string {
	// 매개변수로 받은 byteLen 길이 만큼의 바이트 배열을 생성한다.
	salt := make([]byte, byteLen)

	// io.ReadFull 함수의 첫 번째 매개변수에 rand.Reader을 전달하면 바이트 배열의 인덱스들을 랜덤의 값으로 채울 수 있다.
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		panic(err)
	}

	// 각각의 바이트를 4비트씩 나누어 각각을 16진수로 변환하여 묶은 다음 그것을 그대로 문자열로 변환시켜 반환한다.
	return hex.EncodeToString(salt)
}
