package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

type Keys struct {
	Public  string `mapstructure:"public"`
	Private string `mapstructure:"private"`
}

type KeysSet struct {
	Users map[string]Keys `mapstructure:"users"`
}

func GenerateKeys() (Keys, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return Keys{}, err
	}

	keys := Keys{}
	keys.Public = base64.StdEncoding.EncodeToString(publicKey)
	keys.Private = base64.StdEncoding.EncodeToString(privateKey)

	return keys, nil
}

var ErrKeysNotFound = errors.New("keys not found")
var ErrInvalidKey = errors.New("invalid key")

func (ks KeysSet) ListUsers() []string {
	var usernames = make([]string, 0, len(ks.Users))
	for u, _ := range ks.Users {
		usernames = append(usernames, u)
	}
	return usernames
}

func (ks KeysSet) GetKeys(username string) (Keys, bool) {
	keys, ok := ks.Users[username]
	return keys, ok
}

func (ks KeysSet) GetPrivateKey(username string) (ed25519.PrivateKey, error) {
	keys, ok := ks.Users[username]
	if !ok {
		return ed25519.PrivateKey{}, ErrKeysNotFound
	}

	blob, err := base64.StdEncoding.DecodeString(keys.Private)
	if err != nil {
		return ed25519.PrivateKey{}, ErrInvalidKey
	}

	return blob, nil
}
