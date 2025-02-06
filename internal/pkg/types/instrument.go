package types

import (
	"fmt"
	"strings"
)

type Instrument string

func (i Instrument) Base() string {
	return strings.Split(i.Name(), "/")[0]
}

func (i Instrument) Quote() string {
	return strings.Split(string(i), "/")[1]
}

func (i Instrument) Exchange() string {
	return strings.Split(string(i), ":")[0]
}

func (i Instrument) Name() string {
	return strings.Split(string(i), ":")[1]
}

func (i Instrument) FullName() string {
	return string(i)
}

func (i Instrument) Invert() Instrument {
	return Instrument(fmt.Sprintf("%s:%s/%s", i.Exchange(), i.Base(), i.Quote()))
}

func NewInstrument(base, quote Symbol, exchange string) Instrument {
	return Instrument(fmt.Sprintf("%s:%s/%s", exchange, base, quote))
}
