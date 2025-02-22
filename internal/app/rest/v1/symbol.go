package v1

import (
	"context"
	"database/sql"
	"log/slog"
	"net/url"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"github.com/uptrace/bun"
)

type symbolInput struct {
	// Should be kept at 2x max length of token
	Name string   `path:"name" required:"true" minLength:"3" maxLength:"256"`
	_    struct{} `query:"_" cookie:"_" additionalProperties:"false"`
}

type symbolOutput struct {
	Name                 string   `json:"name"`
	FullName             string   `json:"full_name"`
	Description          string   `json:"description"`
	Type                 string   `json:"type"`
	Session              string   `json:"session"`
	Timezone             string   `json:"timezone"`
	Exchange             string   `json:"exchange"`
	Minmov               int      `json:"minmov"`
	Pricescale           int      `json:"pricescale"`
	HasIntraday          bool     `json:"has_intraday"`
	VisiblePlotsSet      string   `json:"visible_plots_set"`
	HasWeeklyAndMonthly  bool     `json:"has_weekly_and_monthly"`
	SupportedResolutions []string `json:"supported_resolutions"`
	VolumePrecision      int      `json:"volume_precision"`
	DataStatus           string   `json:"data_status"`
}

func SymbolInteractor(ctx context.Context, db *persistence.Database) usecase.IOInteractor {
	supportedResolutions := []string{
		"1S",
		"2S",
		"5S",
		"1",
		"2",
		"5",
		"60",
		"120",
		"300",
		"1D",
		"2D",
		"1W",
		"2W",
		"1M",
		"2M",
		"3M",
	}
	u := usecase.NewInteractor(func(ctx context.Context, input symbolInput, output *symbolOutput) error {
		name, err := url.QueryUnescape(input.Name)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}

		slog.Default().InfoContext(ctx, "Loading instrument", "name", name)
		var m model.Instrument
		if err := db.NewSelect().
			Model(&m).
			Relation("Exchange").
			Where("? = ?", bun.Ident("ins.name"), name).
			Scan(ctx); err != nil {
			if err == sql.ErrNoRows {
				return status.NotFound
			}
			liblog.WithError(ctx, err, "Failed to load instrument.")
			return status.Internal
		}

		output.Name = m.DisplayName
		output.FullName = m.Name
		output.Description = m.Description
		output.Type = cSymbolTypeCrypto
		output.Session = cSymbolDefaultSession
		output.Timezone = cSymbolTimezone
		output.Exchange = m.Exchange.Name
		output.Minmov = cSymbolMinMov
		output.Pricescale = cSymbolPricescale
		output.HasIntraday = true
		output.VisiblePlotsSet = cSymbolVisiblePlotsSets
		output.HasWeeklyAndMonthly = true
		output.SupportedResolutions = supportedResolutions
		output.VolumePrecision = cSymbolVolumePrecision
		output.DataStatus = cSymbolDataStatus
		return nil
	})

	// TODO: Set proper description
	u.SetTitle("Lookup symbol")
	u.SetDescription("Looks up a symbol yo.")

	u.SetExpectedErrors(status.InvalidArgument)

	return u.IOInteractor
}
