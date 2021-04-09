package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// BuildAESKey creates a secret key for encrypting the credentials dependent of the user idm the encrypted corporation name
// and a aes secret key. This permits to use a different key per user and per corporation
func BuildAESKey(userID uint, corporationName, aesSecret string) string {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d:%s:%s", userID, corporationName, aesSecret)))
	return hex.EncodeToString(hash[:])
}

// EncryptString a string given a key
func EncryptString(stringToEncrypt string, aesKey string) (encryptedString string, err error) {
	// Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(aesKey)
	plaintext := []byte(stringToEncrypt)

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	// https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

// DecryptString a string given a key
func DecryptString(encryptedString string, aesKey string) (decryptedString string, err error) {

	key, _ := hex.DecodeString(aesKey)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return
	}

	return fmt.Sprintf("%s", plaintext), nil
}
