package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
)

var screenshotPath = os.Getenv("SCREENSHOTPATH")

func Hash(data string) string {
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

func Path(hash string) string {
	return fmt.Sprintf("%s/%s.jpg", screenshotPath, hash)
}

func Exists(hash string) bool {
	path := fmt.Sprintf("%s/%s.jpg", screenshotPath, hash)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
