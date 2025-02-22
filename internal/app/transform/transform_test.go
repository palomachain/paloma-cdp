package transform

import (
	"math/big"
	"testing"

	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePrice(t *testing.T) {
	tests := []struct {
		name     string
		evt      types.SwapEvent
		i        *model.Instrument
		s0, s1   *model.Symbol
		expected float64
	}{
		{
			name: "Test case 1",
			evt: types.SwapEvent{
				OfferAmount:  big.NewInt(100),
				ReturnAmount: big.NewInt(50),
			},
			i: &model.Instrument{
				Symbol0ID: "symbol0",
				Symbol1ID: "symbol1",
			},
			s0: &model.Symbol{
				ID: "symbol0",
			},
			s1: &model.Symbol{
				ID: "symbol1",
			},
			expected: 2.0,
		},
		{
			name: "Test case 1.1",
			evt: types.SwapEvent{
				OfferAmount:  big.NewInt(50),
				ReturnAmount: big.NewInt(100),
			},
			i: &model.Instrument{
				Symbol0ID: "symbol0",
				Symbol1ID: "symbol1",
			},
			s0: &model.Symbol{
				ID: "symbol0",
			},
			s1: &model.Symbol{
				ID: "symbol1",
			},
			expected: 0.5,
		},
		{
			name: "Test case 2",
			evt: types.SwapEvent{
				OfferAmount:  big.NewInt(50),
				ReturnAmount: big.NewInt(100),
			},
			i: &model.Instrument{
				Symbol0ID: "symbol1",
				Symbol1ID: "symbol0",
			},
			s0: &model.Symbol{
				ID: "symbol0",
			},
			s1: &model.Symbol{
				ID: "symbol1",
			},
			expected: 2.0,
		},
		{
			name: "Test case 2.1",
			evt: types.SwapEvent{
				OfferAmount:  big.NewInt(100),
				ReturnAmount: big.NewInt(50),
			},
			i: &model.Instrument{
				Symbol0ID: "symbol1",
				Symbol1ID: "symbol0",
			},
			s0: &model.Symbol{
				ID: "symbol0",
			},
			s1: &model.Symbol{
				ID: "symbol1",
			},
			expected: 0.5,
		},
		{
			name: "Real data 1",
			evt: types.SwapEvent{
				OfferAmount:  big.NewInt(2000000),
				ReturnAmount: big.NewInt(19560878243514),
			},
			i: &model.Instrument{
				Symbol0ID: "factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd",
				Symbol1ID: "factory/paloma1wff3zdnz2ftptghjk6h6m48m6a5wsnepl7dw526xmvkjy225flys3y036n/MTT.0",
			},
			s0: &model.Symbol{
				ID: "factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd",
			},
			s1: &model.Symbol{
				ID: "factory/paloma1wff3zdnz2ftptghjk6h6m48m6a5wsnepl7dw526xmvkjy225flys3y036n/MTT.0",
			},
			expected: float64(2000000) / float64(19560878243514),
		},
		{
			name: "Real data 2",
			evt: types.SwapEvent{
				OfferAmount:  big.NewInt(19560878243514),
				ReturnAmount: big.NewInt(2000000),
			},
			i: &model.Instrument{
				Symbol0ID: "factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd",
				Symbol1ID: "factory/paloma1wff3zdnz2ftptghjk6h6m48m6a5wsnepl7dw526xmvkjy225flys3y036n/MTT.0",
			},
			s0: &model.Symbol{
				ID: "factory/paloma1wff3zdnz2ftptghjk6h6m48m6a5wsnepl7dw526xmvkjy225flys3y036n/MTT.0",
			},
			s1: &model.Symbol{
				ID: "factory/paloma19nr9tfyx5a0y5f7d5fv978klya6n5cd6lz6zjn9ydqzd7w43gjwq7nc3xy/upusd",
			},
			expected: float64(2000000) / float64(19560878243514),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			price := calculatePrice(tt.evt, tt.i, tt.s0, tt.s1)
			assert.Equal(t, tt.expected, price)
		})
	}
}
