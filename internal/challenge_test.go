package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"
)

func TestCreateAndVerifyChallenge(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	username := "alice"
	ts := time.Date(2024, 4, 1, 12, 0, 0, 0, time.UTC)

	challenge := CreateChallenge(username, ts, priv)

	ok, err := VerifyChallenge(challenge, username, ts, pub)
	if err != nil {
		t.Fatalf("verification error: %v", err)
	}
	if !ok {
		t.Fatal("expected challenge to verify")
	}

	t.Run("altered username", func(t *testing.T) {
		ok, err := VerifyChallenge(challenge, "bob", ts, pub)
		if err != nil {
			t.Fatalf("verification error: %v", err)
		}
		if ok {
			t.Error("expected verification to fail when username changed")
		}
	})

	t.Run("altered timestamp", func(t *testing.T) {
		ok, err := VerifyChallenge(challenge, username, ts.Add(time.Minute), pub)
		if err != nil {
			t.Fatalf("verification error: %v", err)
		}
		if ok {
			t.Error("expected verification to fail when timestamp changed")
		}
	})
}
