package v1

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type barsInput struct {
	// Should be kept at 2x max length of token
	SymbolName string   `path:"name" minLength:"6" maxLength:"256"`
	Resolution string   `query:"resolution"`
	Gte        int      `query:"gte"`
	Lt         int      `query:"lt"`
	_          struct{} `query:"_" cookie:"_" additionalProperties:"false"`
}

type bar struct {
	Time   int     `json:"time"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Volume float64 `json:"volume"`
}
type barsOutput struct {
	Bars []bar `json:"bars"`
}

func BarsInteractor() usecase.IOInteractor {
	// Create use case interactor with references to input/output types and interaction function.
	u := usecase.NewInteractor(func(ctx context.Context, input barsInput, output *barsOutput) error {
		// TODO: Fill
		output.Bars = []bar{
			{
				Time:   1234567890,
				Close:  123.45,
				High:   123.45,
				Low:    123.45,
				Open:   123.45,
				Volume: 123.45,
			},
		}
		return nil
	})

	// Describe use case interactor.
	// TODO: Fill
	u.SetTitle("Search Symbols")
	u.SetDescription("Search for a few symbols.")

	u.SetExpectedErrors(status.InvalidArgument)

	return u.IOInteractor
}
