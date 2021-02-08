package hashing

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const input = "parkjinhong"

func TestGenerateAndCompareSaltedAndPepper(t *testing.T) {
	h := New(generatePeppers(12))

	saltConf := SaltConfig{
		with:    true,
		byteLen: 12,
	}
	pepperConf := PepperConf{
		with: true,
	}

	hash, salt := h.GenerateHash(input, saltConf, pepperConf)

	assert.Equal(t, true, h.Compare(input, salt, hash, pepperConf), "They must same.")
}

func generatePeppers(len int) []string {
	var peppers []string

	for i:=0; i<len; i++ {
		peppers = append(peppers, generateRandomSalt(12))
	}

	return peppers
}