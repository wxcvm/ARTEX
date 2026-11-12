package db

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strings"
)

// Model API keys are optional at-rest secrets. The encryption key is separate
// from PostgreSQL and JWT signing keys; never include it in configuration APIs.
const profileSecretEnv = "ARTEX_LLM_SECRET_KEY"
const profileSecretPrefix = "artex:llm:v1:"
const profileSecretAAD = "artex/llm_profiles/api_key/v1"

var errProfileSecretUnavailable = errors.New("model API key cannot be decrypted; check ARTEX_LLM_SECRET_KEY")
var errProfileSecretConfig = errors.New("ARTEX_LLM_SECRET_KEY must be 64 hexadecimal digits")
var errProfileSecretDowngrade = errors.New("encrypted model API key cannot be replaced without ARTEX_LLM_SECRET_KEY")

func profileSecretCipher() (cipher.AEAD, error) {
	value := strings.TrimSpace(os.Getenv(profileSecretEnv))
	if value == "" {
		return nil, nil
	}
	key, err := hex.DecodeString(value)
	if err != nil || len(key) != 32 {
		return nil, errProfileSecretConfig
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, errProfileSecretConfig
	}
	return cipher.NewGCM(block)
}

// encodeProfileSecret returns exactly one storage representation. Empty keys
// remain empty and legacy installations without a key keep their old behavior.
func encodeProfileSecret(value string) (plain, encrypted string, err error) {
	if value == "" {
		return "", "", nil
	}
	aead, err := profileSecretCipher()
	if err != nil {
		return "", "", err
	}
	if aead == nil {
		return value, "", nil
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", errors.New("cannot generate model API key encryption nonce")
	}
	data := aead.Seal(nonce, nonce, []byte(value), []byte(profileSecretAAD))
	return "", profileSecretPrefix + base64.RawStdEncoding.EncodeToString(data), nil
}

func decodeProfileSecret(plain, encrypted string) (string, error) {
	if encrypted == "" {
		return plain, nil
	}
	aead, err := profileSecretCipher()
	if err != nil {
		return "", err
	}
	if aead == nil || !strings.HasPrefix(encrypted, profileSecretPrefix) {
		return "", errProfileSecretUnavailable
	}
	data, err := base64.RawStdEncoding.DecodeString(strings.TrimPrefix(encrypted, profileSecretPrefix))
	if err != nil || len(data) < aead.NonceSize()+aead.Overhead() {
		return "", errProfileSecretUnavailable
	}
	value, err := aead.Open(nil, data[:aead.NonceSize()], data[aead.NonceSize():], []byte(profileSecretAAD))
	if err != nil {
		return "", errProfileSecretUnavailable
	}
	return string(value), nil
}
