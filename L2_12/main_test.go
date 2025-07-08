package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGrepRegex(t *testing.T) {
	input := "test@example.com\nnot-an-email\nuser@domain.org"

	t.Run("Email regex", func(t *testing.T) {
		r := testGrep(t, input, `\w+@\w+\.\w+`, nil)
		assertEqual(t, r, ">>test@example.com<<\n>>user@domain.org<<")
	})

	t.Run("Fixed string", func(t *testing.T) {
		r := testGrep(t, input, `example.com`, &config{fixed: true})
		assertEqual(t, r, "test@>>example.com<<")
	})
}

func TestGrepBasic(t *testing.T) {
	input := "orange\napple\nbanana\npineapple"

	t.Run("Simple match", func(t *testing.T) {
		r := testGrep(t, input, "apple", nil)
		assertEqual(t, r, ">>apple<<\npine>>apple<<")
	})

	t.Run("Case insensitive", func(t *testing.T) {
		r := testGrep(t, input, "APPLE", &config{ignoreCase: true})
		assertEqual(t, r, ">>apple<<\npine>>apple<<")
	})
}

func TestGrepContext(t *testing.T) {
	input := "1\n2\n3\n4\n5\n6\n7\n8"

	t.Run("Before context", func(t *testing.T) {
		r := testGrep(t, input, "4", &config{before: 2})
		assertEqual(t, r, "2\n3\n>>4<<")
	})

	t.Run("After context", func(t *testing.T) {
		r := testGrep(t, input, "4", &config{after: 2})
		assertEqual(t, r, ">>4<<\n5\n6")
	})
}

func TestGrepCombinedFlags(t *testing.T) {
	input := "ERROR: first\nWARNING: test\nERROR: second\nINFO: message"

	t.Run("Count with inverse", func(t *testing.T) {
		r := testGrep(t, input, "ERROR", &config{count: true, invert: true})
		assertEqual(t, r, "2")
	})

	t.Run("Line numbers with context", func(t *testing.T) {
		r := testGrep(t, input, "WARNING", &config{lineNum: true, around: 1})
		assertEqual(t, r, "2:>>WARNING<<: test")
	})
}

func testGrep(t *testing.T, input, pattern string, cfg *config) string {
	if cfg == nil {
		cfg = &config{}
	}
	cfg.pattern = pattern

	var buf bytes.Buffer
	err := grep(strings.NewReader(input), &buf, *cfg)
	if err != nil {
		t.Fatal(err)
	}

	return strings.TrimSpace(buf.String())
}

func assertEqual(t *testing.T, got, want string) {
	if got != want {
		t.Errorf("\nGot:\n%s\nWant:\n%s", got, want)
	}
}
