package hasher

import (
	"bytes"
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

type Settings struct {
	Times   uint32
	Memory  uint32
	Threads uint8
	KeyLen  uint32
	SaltLen int
}

var (
	passwordSettings = &Settings{
		Times:   1,
		Memory:  1 * 1024,
		Threads: 1,
		KeyLen:  32,
		SaltLen: 4,
	}
)

func generateHashFromSecret(secret string, settings *Settings) ([]byte, error) {
	salt := make([]byte, settings.SaltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}

	hashedPassword := argon2.IDKey([]byte(secret), salt, settings.Times, settings.Memory,
		settings.Threads, settings.KeyLen)
	return append(salt, hashedPassword...), nil
}

func compareHashAndSecret(hash []byte, secret string, settings *Settings) bool {
	salt := hash[0:settings.SaltLen]
	hashedPassword := argon2.IDKey([]byte(secret), salt, settings.Times, settings.Memory,
		settings.Threads, settings.KeyLen)
	return bytes.Equal(hashedPassword, hash[settings.SaltLen:])
}

func GenerateHashFromPassword(password string) ([]byte, error) {
	return generateHashFromSecret(password, passwordSettings)
}

func CompareHashAndPassword(hash []byte, password string) bool {
	return compareHashAndSecret(hash, password, passwordSettings)
}
