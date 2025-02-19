package v1

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type barsInput struct {
	// Should be kept at 2x max length of token
	SymbolName string   `path:"name" minLength:"6" maxLength:"256" required:"true"`
	Resolution string   `query:"resolution" required:"true" enum:"1S,2S,5S,1,2,5,60,120,300,1D,2D,1W,2W,1M,2M,3M"`
	Gte        int64    `query:"gte" required:"true" minimum:"0"`
	Lt         int64    `query:"lt" required:"true" minimum:"0"`
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

func BarsInteractor(db *persistence.Database) usecase.IOInteractor {
	timeBcketMapping := map[string]string{
		"1S":  "1 second",
		"2S":  "2 second",
		"5S":  "5 second",
		"1":   "1 minute",
		"2":   "2 minute",
		"5":   "5 minute",
		"60":  "1 hour",
		"120": "2 hour",
		"300": "5 hour",
		"1D":  "1 day",
		"2D":  "2 day",
		"1W":  "1 week",
		"2W":  "2 week",
		"1M":  "1 month",
		"2M":  "2 month",
		"3M":  "3 month",
	}
	// SELECT
	//   exchange_id,
	//   time_bucket('1 second', time) AS bucket,
	//   min(price) AS low,
	//   max(price) as high,
	//   first(price,time) as open,
	//   last(price,time) as close
	// FROM price_data
	// WHERE
	//   symbol_id=16
	//   AND time > '2025-01-23'
	//   AND time <= '2025-01-25'
	// GROUP BY exchange_id,bucket
	// ORDER BY exchange_id ASC, bucket ASC;
	u := usecase.NewInteractor(func(ctx context.Context, input barsInput, output *barsOutput) error {
		name, err := url.QueryUnescape(input.SymbolName)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		gte := time.Unix(input.Gte, 0)
		lt := time.Unix(input.Lt, 0)
		resolution, ok := timeBcketMapping[input.Resolution]
		if !ok {
			return status.Wrap(errors.New("invalid resolution"), status.InvalidArgument)
		}
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

	// TODO: Fill
	u.SetTitle("Search Symbols")
	u.SetDescription("Search for a few symbols.")

	u.SetExpectedErrors(status.InvalidArgument)

	return u.IOInteractor
}
