package types

import (
	"testing"
)

func TestFromTokenDenom(t *testing.T) {
	tests := []struct {
		input    string
		expected Symbol
		err      bool
	}{
		{"factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd", "UPUSD-19nr9t", false},
		{"factory/paloma1abcd/upusd", "", true},
		{"factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/", "", true},
		{"factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd123", "UPUSD123-19nr9t", false},
		{"factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/up_usd", "UPUSD-19nr9t", false},
	}

	for _, test := range tests {
		result, err := FromTokenDenom(test.input)
		if (err != nil) != test.err {
			t.Errorf("FromTokenDenom(%q) error = %v, expected error = %v", test.input, err, test.err)
		}
		if result != test.expected {
			t.Errorf("FromTokenDenom(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}
