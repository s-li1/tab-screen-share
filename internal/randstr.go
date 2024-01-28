package internal

import (
	"crypto/rand"
  "encoding/hex"
	"fmt"
)

func GenerateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("error:", err)
    return "", err
	}

  return hex.EncodeToString(bytes), nil
}
