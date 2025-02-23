package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type barsInput struct {
	// Should be kept at 2x max length of token
	SymbolName string   `path:"name" minLength:"6" maxLength:"256" required:"true" pattern:"^(DEX|BONDING):[\\S]{3,44}-[[:alnum:]]{6}\\/[\\S]{3,44}-[[:alnum:]]{6}$"`
	Resolution string   `query:"resolution" required:"true" enum:"1S,2S,5S,1,2,5,60,120,300,1D,2D,1W,2W,1M,2M,3M"`
	Gte        int64    `query:"gte" required:"true" minimum:"0"`
	Lt         int64    `query:"lt" required:"true" minimum:"0"`
	_          struct{} `query:"_" cookie:"_" additionalProperties:"false"`
}

type barsOutput struct {
	Bars []model.Bar `json:"bars"`
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

		stmt := db.NewRaw(
			`SELECT
        time_bucket(?, p.time) AS bucket,
        min(p.price) AS low,
        max(p.price) as high,
        first(p.price,p.time) as open,
        last(p.price,p.time) as close
      FROM price_data p JOIN instruments i on p.instrument_id=i.id
      WHERE
        i.name=?
        AND p.time > ?
        AND p.time <= ?
	    GROUP BY bucket
      ORDER BY bucket ASC
      `,
			resolution,
			name,
			gte,
			lt,
		)

		var bars []model.Bar
		err = stmt.Scan(ctx, &bars)
		if err != nil {
			if err == sql.ErrNoRows {
				return status.Wrap(fmt.Errorf("unknown instrument"), status.NotFound)
			}
			liblog.WithError(ctx, err, "Failed to scan bars.")
			return status.Internal
		}
		if bars == nil {
			return status.Wrap(fmt.Errorf("unknown instrument"), status.NotFound)
		}

		output.Bars = bars
		return nil
	})

	// TODO: Fill
	u.SetTitle("Search Symbols")
	u.SetDescription("Search for a few symbols.")

	u.SetExpectedErrors(status.InvalidArgument)

	return u.IOInteractor
}
