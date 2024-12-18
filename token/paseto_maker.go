package token

import (
	"fmt"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
	"time"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey string
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", nil
	}

	token, err := maker.paseto.Encrypt([]byte(maker.symmetricKey), payload, nil)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, []byte(maker.symmetricKey), payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	return payload, err
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be %d length", chacha20poly1305.KeySize)
	}

	pasetoMaker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: symmetricKey,
	}

	return pasetoMaker, nil
}
