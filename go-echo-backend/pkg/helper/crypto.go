package helper

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GetHmac256Hash(text string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(text))

	return hex.EncodeToString(h.Sum(nil))

}
