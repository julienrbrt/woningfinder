package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/woningfinder/woningfinder/internal/util"
)

func TestConfig_Encrypt_Decrypt(t *testing.T) {
	a := assert.New(t)

	key := util.AESGenerateKey()
	stringToEncrypt := "password"

	value, err := util.AESEncrypt(stringToEncrypt, key)
	a.NoError(err)
	a.NotEqual(stringToEncrypt, value)
	value, err = util.AESDecrypt(value, key)
	a.NoError(err)
	a.Equal(stringToEncrypt, value)
}
