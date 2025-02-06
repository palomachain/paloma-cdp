package v1

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type symbolsInput struct {
	UserInput  string   `query:"input" minLength:"3" maxLength:"128"`
	Exchange   *string  `query:"exchange" enum:"palomadex,curvebond"`
	SymbolType *string  `query:"type" pattern:"^crypto$" enum:"crypto"`
	_          struct{} `query:"_" cookie:"_" additionalProperties:"false"`
}

type symbol struct {
	Description string `json:"description"`
	Exchange    string `json:"exchange"`
	Symbol      string `json:"symbol"`
	Ticker      string `json:"ticker"`
	Type        string `json:"type"`
}
type symbolsOutput struct {
	Symbols []symbol `json:"symbols"`
}

func SymbolsInteractor() usecase.IOInteractor {
	// Create use case interactor with references to input/output types and interaction function.
	u := usecase.NewInteractor(func(ctx context.Context, input symbolsInput, output *symbolsOutput) error {
		// TODO: Fill
		output.Symbols = []symbol{
			{
				Description: "A description",
				Exchange:    "palomadex",
				Symbol:      "BTC",
				Ticker:      "BTC",
				Type:        "crypto",
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
