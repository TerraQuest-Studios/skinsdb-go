package help

import (
	"embed"
)

//go:embed *.md
var Files embed.FS

func IsValidPage(page string) bool {
	_, err := Files.ReadFile(page + ".md")
	return err == nil
}
