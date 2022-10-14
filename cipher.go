package newebpay

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

type Cipher struct {
	hashKey []byte `validate:"required,gt=0"`
	hashIV  []byte `validate:"required,gt=0"`

	block cipher.Block
}

func MakeCipher(hashKey, hashIV string) (Cipher, error) {
	c := Cipher{
		hashKey: []byte(hashKey),
		hashIV:  []byte(hashIV),
	}
	err := validate.Struct(c)
	if err != nil {
		return Cipher{}, err
	}
	block, err := aes.NewCipher(c.hashKey)
	if err != nil {
		return Cipher{}, err
	}
	c.block = block
	return c, nil
}

func (c Cipher) encrypt(data []byte) []byte {
	blockMode := cipher.NewCBCEncrypter(c.block, c.hashIV)
	paddedData := pkcs7Pad(data, blockMode.BlockSize())
	encrypted := make([]byte, len(paddedData))
	blockMode.CryptBlocks(encrypted, paddedData)
	return encrypted
}

func (c Cipher) decrypt(data []byte) []byte {
	blockMode := cipher.NewCBCEncrypter(c.block, c.hashIV)
	decrypted := make([]byte, len(data))
	blockMode.CryptBlocks(decrypted, data)
	trimmedData := pkcs7Trim(decrypted)
	return trimmedData
}

func (c Cipher) Encrypt(value string) string {
	data := []byte(value)
	encrypted := c.encrypt(data)
	return hex.EncodeToString(encrypted)
}

func (c Cipher) Decrypt(value string) (string, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}
	decrypted := c.decrypt(data)
	return string(decrypted), nil
}

func (c Cipher) Hash(value string) string {
	encrypted := c.Encrypt(value)
	data := fmt.Sprintf("HashKey=%s&%s&HashIV=%s", string(c.hashKey), encrypted, string(c.hashIV))
	hashed := sha256.Sum256([]byte(data))
	return strings.ToUpper(hex.EncodeToString(hashed[:]))
}
