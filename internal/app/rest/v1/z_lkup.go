package v1

import (
	"context"

	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
)

const (
	cTagAdvancedCharts string = "Advanced Charts"
)

var timeBcketMapping = map[string]string{
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

type lkup struct {
	byName map[string]int64
	byID   map[int64]string
}

func buildExchangeLkup(ctx context.Context, db *persistence.Database) (*lkup, error) {
	lk := &lkup{
		byName: make(map[string]int64),
		byID:   make(map[int64]string),
	}
	var exchanges []model.Exchange
	if err := db.NewSelect().Model(&exchanges).Scan(ctx); err != nil {
		return nil, err
	}

	for _, v := range exchanges {
		lk.byName[v.Name] = v.ID
		lk.byID[v.ID] = v.Name
	}

	return lk, nil
}
