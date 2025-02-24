package v1

import (
	"context"

	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
)

const (
	cTagAdvancedCharts string = "Advanced Charts"
)

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
