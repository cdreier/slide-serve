package main

import (
	"crypto/md5"
	"encoding/hex"
)

func md5Hash(in string) string {
	hasher := md5.New()
	hasher.Write([]byte(in))
	return hex.EncodeToString(hasher.Sum(nil))
}
