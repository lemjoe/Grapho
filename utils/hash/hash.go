package hash

import (
	"crypto/md5"
	"encoding/hex"
)

func GetHash(data []byte) string {
	hash := md5.Sum(data)
	fileName := hex.EncodeToString(hash[:])
	return fileName
}
