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

func TestURLRejectsSchemeless(t *testing.T) {
	err := URL("example.com")
	if err == nil || !strings.HasPrefix(err.Error(), "invalid url") || !strings.Contains(err.Error(), "http or https") {
		t.Fatalf("got %v", err)
	}
}

func TestURLSchemes(t *testing.T) {
	for _, tc := range []struct {
		in      string
		wantErr bool
	}{
		{"ftp://example.com", true},
		{"mailto:a@b.c", true},
		{"HTTP://example.com", false},
	} {
		if err := URL(tc.in); (err != nil) != tc.wantErr {
			t.Errorf("URL(%q) = %v, wantErr %v", tc.in, err, tc.wantErr)
		}
	}
}

func TestURLRejectsTooLong(t *testing.T) {
	s := "https://example.com/" + strings.Repeat("a", 2049-len("https://example.com/"))
	if err := URL(s); err == nil || !strings.HasPrefix(err.Error(), "invalid url") {
		t.Fatalf("got %v", err)
	}
}

func TestURLAcceptsMaxLength(t *testing.T) {
	s := "https://example.com/" + strings.Repeat("a", 2048-len("https://example.com/"))
	if err := URL(s); err != nil {
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
