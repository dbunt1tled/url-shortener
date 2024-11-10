package common

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
)

func isVarType(value interface{}, targetType reflect.Type) bool {
	return reflect.TypeOf(value) == targetType
}

func IsSliceVarOfType(slice interface{}, elemType reflect.Type) bool {
	t := reflect.TypeOf(slice)
	if t.Kind() != reflect.Slice {
		return false
	}
	return t.Elem() == elemType
}

func RandBytes(n uint) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetHexString(bytes []byte) string {
	key := new(big.Int).SetBytes(bytes)
	base16str := fmt.Sprintf("%X", key)
	return base16str
}

func RandStringBytes(length uint) (string, error) {
	var realLength uint
	if length%2 == 0 {
		realLength = length / 2
	} else {
		realLength = length/2 + 1
	}
	bytes, err := RandBytes(realLength)
	if err != nil {
		return "", err
	}
	return GetHexString(bytes)[:length], nil
}
