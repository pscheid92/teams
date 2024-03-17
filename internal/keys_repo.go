package internal

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
)

type KeysRepository interface {
	LoadPrivateKey(username string) (ed25519.PrivateKey, error)

	LoadEncodedPublicKey(username string) (string, error)
	LoadEncodedPrivateKey(username string) (string, error)

	GenerateKeys(username string) (public string, private string, err error)
	DeleteKeys(username string) error
	ListUsers() ([]string, error)
}

type FSKeysRepository struct {
	directory string
}

func NewFSKeyRepository(directory string) *FSKeysRepository {
	return &FSKeysRepository{directory: directory}
}

func (r *FSKeysRepository) LoadPrivateKey(username string) (ed25519.PrivateKey, error) {
	encoded, err := r.LoadEncodedPrivateKey(username)
	if err != nil {
		return nil, err
	}

	return base64.StdEncoding.DecodeString(encoded)
}

func (r *FSKeysRepository) LoadEncodedPublicKey(username string) (string, error) {
	path := r.publicPath(username)
	blob, err := os.ReadFile(path)
	return string(blob), err
}

func (r *FSKeysRepository) LoadEncodedPrivateKey(username string) (string, error) {
	path := r.privatePath(username)
	blob, err := os.ReadFile(path)
	return string(blob), err
}

func (r *FSKeysRepository) GenerateKeys(username string) (public string, private string, err error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", "", err
	}

	publicBlob := make([]byte, base64.StdEncoding.EncodedLen(len(publicKey)))
	base64.StdEncoding.Encode(publicBlob, publicKey)
	if err := os.WriteFile(r.publicPath(username), publicBlob, os.ModePerm); err != nil {
		return "", "", err
	}

	privateBlob := make([]byte, base64.StdEncoding.EncodedLen(len(privateKey)))
	base64.StdEncoding.Encode(privateBlob, privateKey)
	if err := os.WriteFile(r.privatePath(username), privateBlob, os.ModePerm); err != nil {
		return "", "", err
	}

	return string(publicBlob), string(privateBlob), nil
}

func (r *FSKeysRepository) DeleteKeys(username string) error {
	if err := os.Remove(r.publicPath(username)); err != nil {
		return err
	}

	if err := os.Remove(r.privatePath(username)); err != nil {
		return err
	}

	return nil
}

func (r *FSKeysRepository) ListUsers() ([]string, error) {
	files, err := os.ReadDir(r.directory)
	if err != nil {
		return nil, err
	}

	usersSet := make(map[string]struct{})
	for _, file := range files {
		filename := file.Name()
		username := strings.TrimSuffix(strings.TrimSuffix(filename, ".public"), ".private")
		usersSet[username] = struct{}{}
	}

	users := make([]string, 0, len(usersSet))
	for user := range usersSet {
		users = append(users, user)
	}

	return users, nil
}

func (r *FSKeysRepository) publicPath(username string) string {
	return filepath.Join(r.directory, username+".public")
}

func (r *FSKeysRepository) privatePath(username string) string {
	return filepath.Join(r.directory, username+".private")
}
