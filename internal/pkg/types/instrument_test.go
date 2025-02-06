package types

import (
	"testing"
)

func TestInstrument(t *testing.T) {
	tests := []struct {
		exchange string
		base     string
		quote    string
		expected string
		name     string
	}{
		{"binance", "BTC", "USDT", "binance:BTC/USDT", "BTC/USDT"},
		{"kraken", "ETH", "EUR", "kraken:ETH/EUR", "ETH/EUR"},
		{"coinbase", "LTC", "USD", "coinbase:LTC/USD", "LTC/USD"},
		{"palomadex", "UPUSD-19nr9t", "MTT.0-1wff3z", "palomadex:UPUSD-19nr9t/MTT.0-1wff3z", "UPUSD-19nr9t/MTT.0-1wff3z"},
	}

	for _, test := range tests {
		i := NewInstrument(Symbol(test.base), Symbol(test.quote), test.exchange)
		fn := i.FullName()
		if fn != test.expected {
			t.Errorf("Instrument.Fullname() = %v, expected %v", fn, test.expected)
		}
		n := i.Name()
		if n != test.name {
			t.Errorf("Instrument.Fullname() = %v, expected %v", n, test.expected)
		}
		if i.Exchange() != test.exchange {
			t.Errorf("Instrument.Exchange() = %v, expected %v", i.Exchange(), test.exchange)
		}
		if i.Base() != test.base {
			t.Errorf("Instrument.Base() = %v, expected %v", i.Base(), test.base)
		}
		if i.Quote() != test.quote {
			t.Errorf("Instrument.Quote() = %v, expected %v", i.Quote(), test.quote)
		}
	}
}
