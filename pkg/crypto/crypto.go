package crypto

import (
	"crypto/md5"
	"fmt"
	"math/rand"
)

func GetMD5(data string) string {
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}

func GetToken32() string {
	token := make([]byte, 32)
	rand.Read(token)
	return fmt.Sprintf("%x", token)
}
