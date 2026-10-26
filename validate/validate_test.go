package validate

import (
	"strings"
	"testing"
)

func TestURLAcceptsAbsolute(t *testing.T) {
	if err := URL("https://example.com/a?b=1#c"); err != nil {
		t.Fatal(err)
	}
}

func TestURLAcceptsSchemeless(t *testing.T) {
	if err := URL("example.com"); err != nil {
		t.Fatal(err)
	}
}

func TestURLAcceptsLong(t *testing.T) {
	if err := URL("https://example.com/" + strings.Repeat("a", 5000)); err != nil {
		t.Fatal(err)
	}
}

func TestURLRejectsEmpty(t *testing.T) {
	if err := URL(""); err == nil || !strings.HasPrefix(err.Error(), "invalid url") {
		t.Fatalf("got %v", err)
	}
}

func TestURLRejectsUnparsable(t *testing.T) {
	if err := URL("http://exa mple.com"); err == nil {
		t.Fatal("expected error")
	}
}
