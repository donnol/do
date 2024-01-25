package do

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
)

type Signer struct {
	secret string // base64 encoded
}

// NewSigner which secret is base64 encoded
//
// secret can use `ssh-keygen -t rsa -b 2048 -m PKCS8 -f jwtRS256.key` to generate, it will generate `jwtRS256.key` and `jwtRS256.key.pub`, and the first file is what we want
func NewSigner(secret string) *Signer {
	return &Signer{
		secret: secret,
	}
}

// Sign return raw's signature which is base64 encoded
func (s *Signer) Sign(raw []byte) (r string, err error) {
	secret, err := base64.StdEncoding.DecodeString(s.secret)
	if err != nil {
		return
	}

	key, err := x509.ParsePKCS8PrivateKey(secret)
	if err != nil {
		return
	}

	hashed := sha256.Sum256(raw)
	enc, err := rsa.SignPKCS1v15(nil, key.(*rsa.PrivateKey), crypto.SHA256, hashed[:])
	if err != nil {
		return
	}

	r = base64.StdEncoding.EncodeToString(enc)

	return
}
