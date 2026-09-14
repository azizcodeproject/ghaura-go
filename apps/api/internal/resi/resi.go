package resi

import (
	"crypto/rand"
	"fmt"
	"time"
)

const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func Generate() (string, error) {
	datePart := time.Now().UTC().Format("060102")
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("generate resi suffix: %w", err)
	}

	suffix := make([]byte, 6)
	for index := range suffix {
		suffix[index] = alphabet[int(randomBytes[index])%len(alphabet)]
	}

	return fmt.Sprintf("GHR-%s-%s", datePart, string(suffix)), nil
}
