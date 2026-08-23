package main

import (
	"strings"
	"testing"
)

func TestHighlightRedWrapsText(t *testing.T) {
	got := highlightRed("hello")
	if !strings.Contains(got, "hello") {
		t.Fatalf("highlightRed missing original text: %q", got)
	}
	if !strings.HasPrefix(got, "\x1b[") {
		t.Fatalf("highlightRed missing ANSI prefix: %q", got)
	}
	if !strings.HasSuffix(got, "\x1b[0m") {
		t.Fatalf("highlightRed missing reset suffix: %q", got)
	}
}
