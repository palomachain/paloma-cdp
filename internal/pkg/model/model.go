package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Exchange struct {
	bun.BaseModel `bun:"table:exchanges,alias:e"`

	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",notnull"`
}

type Symbol struct {
	bun.BaseModel `bun:"table:symbols,alias:s"`

	ID          int64     `bun:",pk,autoincrement"`
	Name        string    `bun:",notnull"`
	Description string    `bun:",notnull"`
	ExchangeID  int64     `bun:",notnull"`
	Exchange    *Exchange `bun:"rel:belongs-to,join:exchange_id=id"`
}

type PriceData struct {
	bun.BaseModel `bun:"table:price_data,alias:p"`

	SymbolID int64     `bun:",notnull"`
	Symbol   *Symbol   `bun:"rel:belongs-to,join:symbol_id=id"`
	Time     time.Time `bun:",notnull"`
	Price    float64   `bun:",notnull"`
}
