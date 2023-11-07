package do

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"strings"
)

// https://www.kancloud.cn/wizardforcel/golang-stdlib-ref/121494
// https://csrc.nist.gov/projects/block-cipher-techniques/bcm/current-modes
// https://www.cnblogs.com/happyhippy/archive/2006/12/23/601353.html

type Crypto struct {
	key string
}

// NewCrypto NewCrypto
// key: 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func NewCrypto(key string) (*Crypto, error) {
	var ik string

	keyLen := len(key)
	switch keyLen {
	case 16, 24, 32:
		ik = key
	default:
		return nil, fmt.Errorf("bad key length: %d", keyLen)
	}

	return &Crypto{
		key: ik,
	}, nil
}

func (e Crypto) Encrypt(money string) (r string, err error) {
	key := []byte(e.key)
	plaintext := []byte(money)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// 生成随机字节
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	r = string(ciphertext)

	return
}

func (e Crypto) Decrypt(r string) (money string, err error) {
	if strings.TrimSpace(r) == "" {
		return
	}

	key := []byte(e.key)
	ciphertext := []byte(r)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// 取出随机字节
	if len(ciphertext) < aes.BlockSize {
		err = fmt.Errorf("ciphertext too short: %s", r)
		return
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	money = string(ciphertext)

	return
}
