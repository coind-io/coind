package crypto

import (
	"crypto/rand"

	"golang.org/x/crypto/ed25519"
)

func GenerateKeyPair() (*PrivKey, *PubKey, error) {
	edpub, edkey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	privkey, err := NewPrivKeyFromBytes(edkey.Seed())
	if err != nil {
		return nil, nil, err
	}
	pubkey, err := NewPubKeyFromBytes([]byte(edpub))
	if err != nil {
		return nil, nil, err
	}
	return privkey, pubkey, nil
}
