package lib

import (
	"crypto/md5"
	"fmt"
)

// CalculateMD5 calculates MD5 hash by string
func CalculateMD5(str string) string {
	data := []byte(str)
	return fmt.Sprintf("%x", md5.Sum(data))
}
