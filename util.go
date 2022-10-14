package newebpay

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

func pkcs7Pad(data []byte, blockSize int) []byte {
	paddingLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(paddingLen)}, paddingLen)
	return append(data, padding...)
}

func pkcs7Trim(data []byte) []byte {
	dataLen := len(data)
	paddingLen := int(data[dataLen-1])
	return data[:(dataLen - paddingLen)]
}

func sha256HashHex(value string) string {
	hashed := sha256.Sum256([]byte(value))
	return strings.ToUpper(hex.EncodeToString(hashed[:]))
}

var validate = newValidator()

func splitValidator(fl validator.FieldLevel) bool {
	acceptedValues := strings.Split(fl.Param(), " ")
	field := fl.Field()
	var reflected string
	switch field.Kind() {
	case reflect.String:
		reflected = field.String()
	default:
		panic(fmt.Sprintf("Bad field type %T", field.Interface()))
	}
	values := strings.Split(reflected, ",")
	valueMap := make(map[string]bool)
	for _, acceptedValue := range acceptedValues {
		valueMap[acceptedValue] = false
	}
	for _, value := range values {
		if v, ok := valueMap[value]; !ok || v {
			return false
		}
		valueMap[value] = true
	}
	return true
}

func newValidator() *validator.Validate {
	validator := validator.New()
	validator.RegisterValidation("split", splitValidator)
	return validator
}
