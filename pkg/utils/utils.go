package utils

import "crypto/sha1"

//Get SHA1 of string
func Sha1(s string) string {
	h := sha1.New()

	h.Write([]byte(s))
	bs := h.Sum(nil)

	hash := string(bs)
	return hash
}
