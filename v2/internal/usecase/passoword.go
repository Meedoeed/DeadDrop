package usecase

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const (
	lower  = "abcdefghijklmnopqrstuvwxyz"
	upper  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits = "0123456789"
	spec   = "@!$%&#"
	total  = lower + upper + digits + spec
)

func GeneratePassword(length int) (string, error) {
	if length < 5 {
		return "", fmt.Errorf("length of password is too low")
	}

	password := make([]byte, length)

	randChar := func(set string) (byte, error) {
		max := big.NewInt(int64(len(set)))
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return 0, err
		}
		return set[n.Int64()], nil
	}

	chars := []string{lower, upper, spec, digits}
	for i := 0; i < 4; i++ {
		char, err := randChar(chars[i])
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	for i := 4; i < length; i++ {
		char, err := randChar(total)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	for i := length - 1; i > 0; i-- {
		max := big.NewInt(int64(i + 1))
		j, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		password[i], password[j.Int64()] = password[j.Int64()], password[i]
	}

	return string(password), nil
}
