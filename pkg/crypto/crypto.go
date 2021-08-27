package crypto

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(data string) string {
	hash := md5.Sum([]byte(data))
	return fmt.Sprintf("%x", hash)
}
