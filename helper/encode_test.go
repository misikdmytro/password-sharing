package helper

import (
	"fmt"
	"testing"

	"github.com/misikdmitriy/password-sharing/config"
)

func TestEncodeDecodeShouldDoIt(t *testing.T) {
	config := &config.Config{}
	config.Encrypt.Secret = "123456789123456789012345"
	config.Encrypt.IV = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}

	encoder := NewEncoder(config)
	text := "initial text"

	encoded, err := encoder.Encode(text)
	if err != nil {
		t.Error(err)
	}

	decoded, err := encoder.Decode(encoded)
	if err != nil {
		t.Error(err)
	}

	if text != decoded {
		t.Error(fmt.Errorf("expected decoded to be '%s' but was '%s'", text, decoded))
	}
}
