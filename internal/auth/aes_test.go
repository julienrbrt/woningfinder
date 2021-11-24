package auth_test

import (
	"testing"

	"github.com/julienrbrt/woningfinder/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestConfig_Encrypt_Decrypt(t *testing.T) {
	a := assert.New(t)

	key := auth.BuildAESKey(1, "foo@bar.com", "foo")
	stringToEncrypt := "password"

	value, err := auth.EncryptString(stringToEncrypt, key)
	a.NoError(err)
	a.NotEqual(stringToEncrypt, value)
	value, err = auth.DecryptString(value, key)
	a.NoError(err)
	a.Equal(stringToEncrypt, value)
}
