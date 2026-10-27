package validate

import (
	"errors"
	"fmt"
	"net/url"
	"unicode/utf8"
)

const maxLength = 2048

func URL(s string) error {
	if s == "" {
		return errors.New("invalid url: empty")
	}
	if utf8.RuneCountInString(s) > maxLength {
		return errors.New("invalid url: longer than 2048 characters")
	}
	u, err := url.Parse(s)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid url: scheme must be http or https")
	}
	return nil
}
