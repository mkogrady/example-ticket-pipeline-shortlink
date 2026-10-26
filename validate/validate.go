package validate

import (
	"errors"
	"fmt"
	"net/url"
)

func URL(s string) error {
	if s == "" {
		return errors.New("invalid url: empty")
	}
	if _, err := url.Parse(s); err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}
	return nil
}
