package helpers

import (
	"crypto/sha512"
)

func HashSha512(data string) string {
	return string(sha512.New().Sum([]byte(data)))
}
