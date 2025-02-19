package v1

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"github.com/uptrace/bun"
)

const (
	cSymbolTypeCrypto       = "crypto"
	cSymbolDefaultSession   = "24x7"
	cSymbolTimezone         = "Etc/UTC"
	cSymbolMinMov           = 1
	cSymbolPricescale       = 10000
	cSymbolVisiblePlotsSets = "ohlcv"
	cSymbolVolumePrecision  = 6
	cSymbolDataStatus       = "streaming"
)

type symbolsInput struct {
	UserInput  string   `query:"input" required:"true" minLength:"3" maxLength:"128"`
	Exchange   *string  `query:"exchange" enum:"DEX,BONDING"`
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

func SymbolsInteractor(ctx context.Context, db *persistence.Database) usecase.IOInteractor {
	// We can avoid joins on every query by retrieving static exchange info once
	lkup, err := buildExchangeLkup(ctx, db)
	if err != nil {
		liblog.WithError(ctx, err, "Failed to load exchange lookup.")
		panic(err)
	}

	u := usecase.NewInteractor(func(ctx context.Context, input symbolsInput, output *symbolsOutput) error {
		userInput, err := url.QueryUnescape(input.UserInput)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		var m []model.Instrument
		stmt := db.NewSelect().Model(&m).Where("? LIKE ?", bun.Ident("name"), fmt.Sprintf("%%%s%%", strings.ToUpper(userInput)))

		if input.Exchange != nil {
			eid, ok := lkup.byName[*input.Exchange]
			if !ok {
				return status.Wrap(errors.New("unknown exchange"), status.InvalidArgument)
			}
			stmt.Where("? = ?", bun.Ident("exchange_id"), eid)
		}
		err = stmt.Scan(ctx)
		if err != nil {
			liblog.WithError(ctx, err, "Failed to load symbols.")
			return status.Internal
		}

		symbols := make([]symbol, len(m))
		for i, v := range m {
			symbols[i] = symbol{
				Description: v.Description,
				Exchange:    lkup.byID[v.ExchangeID],
				Symbol:      v.DisplayName,
				Ticker:      v.ShortName,
				Type:        cSymbolTypeCrypto,
			}
		}

		output.Symbols = symbols
		return nil
	})

	u.SetTitle("Search Symbols")
	u.SetDescription("Search for a few symbols.")

	u.SetExpectedErrors(status.InvalidArgument)

	return u.IOInteractor
}
