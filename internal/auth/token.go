package auth

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

var (
	cryptoKey  = []byte("12345678901234567890123456789012") // 32 bytes for AES-256
	initVector = []byte("1234567890abcdef")                 // 16 bytes for CBC block size
)

// PKCS7Padding pads the input text to be a multiple of the block size (required for CBC mode)
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - (len(ciphertext) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7Unpadding removes the padding after decryption
func PKCS7Unpadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("decryption failed: empty source")
	}
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, fmt.Errorf("decryption failed: invalid padding size")
	}
	return src[:(length - unpadding)], nil
}

// EncryptToken matches Node.js cipher.update(token, "utf8", "hex") + cipher.final("hex")
func EncryptToken(token string) (string, error) {
	block, err := aes.NewCipher(cryptoKey)
	if err != nil {
		return "", err
	}

	// 1. Pad the plaintext to match AES block size (16 bytes)
	paddedPlaintext := PKCS7Padding([]byte(token), aes.BlockSize)

	// 2. Encrypt using CBC mode
	ciphertext := make([]byte, len(paddedPlaintext))
	mode := cipher.NewCBCEncrypter(block, initVector)
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	// 3. Return as hex string
	return hex.EncodeToString(ciphertext), nil
}

// DecryptToken takes the hex string and decrypts it back to original plain text
func DecryptToken(hexCiphertext string) (string, error) {
	// 1. Decode hex string back to bytes
	ciphertext, err := hex.DecodeString(hexCiphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(cryptoKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// 2. Decrypt using CBC mode
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, initVector)
	mode.CryptBlocks(plaintext, ciphertext)

	// 3. Remove padding
	unpaddedPlaintext, err := PKCS7Unpadding(plaintext)
	if err != nil {
		return "", err
	}

	return string(unpaddedPlaintext), nil
}
