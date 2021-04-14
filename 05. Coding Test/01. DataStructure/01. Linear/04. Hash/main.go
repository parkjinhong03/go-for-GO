package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func main() {
	noHashed := []struct{
		k, v string
	}{
		{k: , v: },
	}

	// sha256 알고리즘으로 해싱된 Byte Array 반환
	hash := sha256.Sum256([]byte(""))
	
	// Byte Array를 Slice로 변환한 후 16진수로 인코딩
	enc := hex.EncodeToString(hash[:])
	
	fmt.Println(enc)
}
