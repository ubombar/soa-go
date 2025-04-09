package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/spf13/viper"

	"github.com/ubombar/soa/internal/datetime"
)

var ErrBadFilename = errors.New("given question title is longer than 255 characters")

func SanitizeName(title string) (string, error) {
	if utf8.RuneCountInString(title) > 255 {
		return "", ErrBadFilename
	}

	sanitizedTitle := strings.Map(func(r rune) rune {
		switch r {
		case '/':
			return '-'
		case ':':
			return '-'
		default:
			return r
		}
	}, title)

	return sanitizedTitle, nil
}

func GetFilename(dir string, filename string) string {
	vaultDir := viper.GetString("vault-dir")
	return filepath.Join(vaultDir, dir, filename)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func LiteratureFilename(pdfPath string, date datetime.Date) string {
	baseName := filepath.Base(pdfPath)
	noteName := fmt.Sprintf("L %s %s.md", date.String(), baseName)
	return noteName
}

func QuestionFilename(title string, date datetime.Date) string {
	noteName := fmt.Sprintf("Q %s %s.md", date.String(), title)
	return noteName
}
