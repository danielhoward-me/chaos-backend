package screenshot

import (
	_ "embed"
	"fmt"
	"os"
)

var screenshotPath = os.Getenv("SCREENSHOTPATH")

//go:embed placeholder.jpg
var PlaceholderImage []byte

func Path(hash string) string {
	return fmt.Sprintf("%s/%s.jpg", screenshotPath, hash)
}

func Exists(hash string) bool {
	path := fmt.Sprintf("%s/%s.jpg", screenshotPath, hash)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
