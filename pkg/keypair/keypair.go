package keypair

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

type KeyPair struct {
	Private string `json:"private"`
	Public  string `json:"public"`
}

func Generate() (KeyPair, error) {
	priv, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		return KeyPair{}, fmt.Errorf("generate x25519 key: %w", err)
	}

	return KeyPair{
		Private: base64.StdEncoding.EncodeToString(priv.Bytes()),
		Public:  base64.StdEncoding.EncodeToString(priv.PublicKey().Bytes()),
	}, nil
}

func DerivePublicKey(privateKey string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(strings.TrimSpace(privateKey))
	if err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}

	key, err := ecdh.X25519().NewPrivateKey(b)
	if err != nil {
		return "", fmt.Errorf("invalid key: %w", err)
	}

	return base64.StdEncoding.EncodeToString(key.PublicKey().Bytes()), nil
}
