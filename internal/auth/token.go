package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// secretbox requires exactly a 32-byte key array
var secretKey = [32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2}

func EncryptToken(token string) (string, error) {
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return "", err
	}

	// This single line encrypts, seals, and adds authentication tags without padding boilerplate
	encrypted := secretbox.Seal(nonce[:], []byte(token), &nonce, &secretKey)
	return hex.EncodeToString(encrypted), nil
}

func DecryptToken(hexCiphertext string) (string, error) {
	data, err := hex.DecodeString(hexCiphertext)
	if err != nil {
		return "", err
	}

	if len(data) < 24 {
		return "", fmt.Errorf("ciphertext too short")
	}

	var decryptNonce [24]byte
	copy(decryptNonce[:], data[:24])

	// Decrypts completely behind the scenes
	decrypted, ok := secretbox.Open(nil, data[24:], &decryptNonce, &secretKey)
	if !ok {
		return "", fmt.Errorf("decryption failed or data tampered with")
	}

	return string(decrypted), nil
}
