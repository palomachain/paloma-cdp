package types

import (
	"fmt"
	"regexp"
	"strings"
)

type Symbol string

func (s Symbol) String() string {
	return string(s)
}

// SymbolFromTokenDenom converts a token denomination to a Symbol.
// It expects input in the form of tokenfactory token names,
// i.e. factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd
// It returns a symbol representation of the token name, formed like this:
// SUBDENOM-<First 6 letters of bech32 address - prefix)
// example: UPUSD-19nr9t
// Any special characters in the token name are removed.
func SymbolFromTokenDenom(i string) (Symbol, error) {
	parts := strings.Split(i, "/")
	if len(parts) < 3 {
		return "", fmt.Errorf("invalid token denomination format")
	}
	if len(parts) > 3 {
		parts = append(parts[:2], strings.Join(parts[2:], "/"))
	}

	address := parts[1]
	subdenom := parts[2]

	if len(address) < 20 {
		return "", fmt.Errorf("invalid address length")
	}

	// Extract the first 6 characters of the address
	// omitting the paloma prefix.
	addressPart := address[6:12]

	// Remove any special characters from the subdenom
	re := regexp.MustCompile(`[^a-zA-Z0-9.]+`)
	cleanSubdenom := re.ReplaceAllString(subdenom, "")
	if len(cleanSubdenom) < 3 {
		return "", fmt.Errorf("invalid subdenom")
	}

	// Form the symbol
	symbol := strings.ToUpper(cleanSubdenom) + "-" + addressPart

	return Symbol(symbol), nil
}
