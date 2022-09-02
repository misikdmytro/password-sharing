package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"

	"github.com/misikdmitriy/password-sharing/config"
)

type Encoder interface {
	Encode(data string) (string, error)
	Decode(data string) (string, error)
}

type encoder struct {
	config *config.Config
}

func NewEncoder(config *config.Config) Encoder {
	return &encoder{
		config: config,
	}
}

func (e *encoder) Encode(data string) (string, error) {
	block, err := aes.NewCipher([]byte(e.config.Encrypt.Secret))
	if err != nil {
		return "", err
	}

	plainText := []byte(data)
	cfb := cipher.NewCFBEncrypter(block, e.config.Encrypt.IV)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (e *encoder) Decode(data string) (string, error) {
	block, err := aes.NewCipher([]byte(e.config.Encrypt.Secret))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, e.config.Encrypt.IV)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
