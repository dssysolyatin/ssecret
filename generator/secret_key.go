package generator

import "crypto/rand"

type SecretKeyGenerator interface {
	Generate() ([]byte, error)
}

var _ SecretKeyGenerator = (*randSecretKeyGenerator)(nil)

type randSecretKeyGenerator struct {
}

func NewRandSecretKeyGenerator() *randSecretKeyGenerator {
	return &randSecretKeyGenerator{}
}

func (r *randSecretKeyGenerator) Generate() ([]byte, error) {
	b := make([]byte, 128)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
