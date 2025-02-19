package model

import (
	"time"

	"github.com/cometbft/cometbft/types"
	"github.com/uptrace/bun"
)

type Exchange struct {
	bun.BaseModel `bun:"table:exchanges,alias:e"`

	ID   int64  `bun:",pk,autoincrement"`
	Name string `bun:",notnull"`
}

type ExchangeLkup struct {
	bun.BaseModel `bun:"table:exchange_lkup,alias:el"`

	Address    string    `bun:",pk"`
	ExchangeID int64     `bun:",notnull"`
	Exchange   *Exchange `bun:"rel:belongs-to,join:exchange_id=id"`
}

type Symbol struct {
	bun.BaseModel `bun:"table:symbols,alias:s"`

	ID          string `bun:",pk"`
	DisplayName string `bun:",notnull"`
}

type Instrument struct {
	bun.BaseModel `bun:"table:instruments,alias:ins"`

	ID          int64     `bun:",pk,autoincrement"`
	Name        string    `bun:",notnull"`
	ShortName   string    `bun:",notnull"`
	DisplayName string    `bun:",notnull"`
	Description string    `bun:",notnull"`
	ExchangeID  int64     `bun:",notnull"`
	Exchange    *Exchange `bun:"rel:belongs-to,join:exchange_id=id"`
	Symbol0ID   string    `bun:"symbol0_id,notnull"`
	Symbol0     *Symbol   `bun:"rel:belongs-to,join:symbol0_id=id"`
	Symbol1ID   string    `bun:"symbol1_id,notnull"`
	Symbol1     *Symbol   `bun:"rel:belongs-to,join:symbol1_id=id"`
}

type PriceData struct {
	bun.BaseModel `bun:"table:price_data,alias:p"`

	InstrumentID int64     `bun:",notnull"`
	Instrument   *Symbol   `bun:"rel:belongs-to,join:instrument_id=id"`
	Time         time.Time `bun:",notnull"`
	Price        float64   `bun:",notnull"`
}

type RawTransaction struct {
	bun.BaseModel `bun:"table:ingest,alias:i"`

	Hash      string            `bun:"tx_hash,pk"`
	Data      types.EventDataTx `bun:",notnull,msgpack"`
	Timestamp time.Time         `bun:"received,notnull,default:current_timestamp"`
}
