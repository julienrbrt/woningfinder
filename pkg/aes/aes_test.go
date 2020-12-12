package aes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/pkg/aes"
)

func TestConfig_Encrypt_Decrypt(t *testing.T) {
	a := assert.New(t)

	key := aes.GenerateKey()
	stringToEncrypt := "password"

	value, err := aes.Encrypt(stringToEncrypt, key)
	a.NoError(err)
	a.NotEqual(stringToEncrypt, value)
	value, err = aes.Decrypt(value, key)
	a.NoError(err)
	a.Equal(stringToEncrypt, value)
}
