package codeutil

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateFourDigitsCode() (string, error) {
	var code string

	for i := 0; i < 4; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(4))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n.Int64())
	}

	return code, nil
}
