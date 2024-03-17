package internal

import (
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"time"
)

const (
	MessageSeparator = "__"
	TimeFormat       = time.RFC3339
)

func CreateChallenge(username string, timestamp time.Time, key ed25519.PrivateKey) string {
	message := generateMessage(username, timestamp)
	signature := ed25519.Sign(key, message)
	return base64.StdEncoding.EncodeToString(signature)
}

func VerifyChallenge(challenge string, username string, timestamp time.Time, key ed25519.PublicKey) (bool, error) {
	signature, err := base64.StdEncoding.DecodeString(challenge)
	if err != nil {
		return false, fmt.Errorf("verify challenge: %w", err)
	}

	message := generateMessage(username, timestamp)
	return ed25519.Verify(key, message, signature), nil
}

func generateMessage(username string, timestamp time.Time) []byte {
	return []byte(username + MessageSeparator + timestamp.Format(TimeFormat))
}
