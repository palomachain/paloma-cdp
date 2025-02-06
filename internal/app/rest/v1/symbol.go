package v1

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type symbolInput struct {
	// Should be kept at 2x max length of token
	Name string   `path:"name" minLength:"3" maxLength:"256"`
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

func SymbolInteractor() usecase.IOInteractor {
	// Create use case interactor with references to input/output types and interaction function.
	u := usecase.NewInteractor(func(ctx context.Context, input symbolInput, output *symbolOutput) error {
		// TODO: Fill
		return nil
	})

	// Describe use case interactor.
	// TODO: Fill
	u.SetTitle("Lookup symbol")
	u.SetDescription("Looks up a symbol yo.")

	u.SetExpectedErrors(status.InvalidArgument)

	return u.IOInteractor
}
