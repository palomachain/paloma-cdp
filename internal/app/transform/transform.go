package transform

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/palomachain/paloma-cdp/internal/pkg/types"
	"github.com/uptrace/bun"
)

type Configuration struct {
	PollingInterval time.Duration `env:"CDP_TRANSFORMER_POLLING_INTERVAL,notEmpty" envDefault:"1s"`
}

const (
	cPollingLimit int = 20
)

var gLkUp map[string]model.Exchange

func Run(
	ctx context.Context,
	v string,
	db *persistence.Database,
	cfg *Configuration,
) error {
	slog.Default().InfoContext(ctx, "Service running.", "version", v)

	var err error
	gLkUp, err = loadExchangeLkUp(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to load exchange lookup: %w", err)
	}

	// TODO: When and HOW do we recover from a panic?
	tkr := time.NewTicker(cfg.PollingInterval)
	for {
		select {
		case <-ctx.Done():
			tkr.Stop()
			return nil
		case <-tkr.C:
			if err := handleTick(ctx, db); err != nil {
				liblog.WithError(ctx, err, "Failed to perform data transformation!")
			}
		}
	}
}

func handleTick(ctx context.Context, db *persistence.Database) error {
	// Ensure we keep scanning for results, even if there are more than 20
	// We do not want to wait for the next tick to do this.
	for {
		var txs []model.RawTransaction
		count, err := db.NewSelect().Model(&txs).Limit(cPollingLimit).ScanAndCount(ctx)
		if err != nil {
			return fmt.Errorf("failed to fetch transactions: %w", err)
		}
		if count < 1 {
			return nil
		}

		slog.Default().DebugContext(ctx, "Fetched transactions.", "count", count)
		err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
			for _, v := range txs {
				if err := consumeTx(ctx, tx, gLkUp, v); err != nil {
					return err
				}
			}

			_, err = tx.NewDelete().Model(&txs).WherePK().Exec(ctx)
			return err
		})
		if err != nil {
			return err
		}

		// Preserve any intermiate entered data until next tick.
		if count < cPollingLimit {
			return nil
		}
	}
}

func consumeTx(ctx context.Context,
	db bun.Tx,
	exchangeLkup map[string]model.Exchange,
	tx model.RawTransaction,
) (err error) {
	slog.Default().InfoContext(ctx, "Consuming transaction", "hash", tx.Hash)
	defer func() {
		if err != nil {
			// Move TX to DLQ and continue normal execution
			_, err = db.NewInsert().Model(&tx).ModelTableExpr("ingest_dlq").Exec(ctx)
		}
	}()

	evts, err := types.TryParseSwapEvents(ctx, tx.Data.Result.Events)
	if err != nil {
		return fmt.Errorf("failed to parse swap events: %w", err)
	}

	if len(evts) < 1 {
		slog.Default().InfoContext(ctx, "No swap events found", "hash", tx.Hash)
		return nil
	}

	for _, evt := range evts {
		evt.Timestamp = tx.Timestamp
		if err := consumeEvent(ctx, db, exchangeLkup, evt); err != nil {
			liblog.WithError(ctx, err, "Failed to consume event", "hash", tx.Hash, "event", evt)
			return err
		}
	}

	return nil
}

func consumeEvent(ctx context.Context,
	db bun.Tx,
	exchangeLkup map[string]model.Exchange,
	evt types.SwapEvent,
) error {
	exchange, ok := exchangeLkup[evt.Sender]
	if !ok {
		return fmt.Errorf("sender %s not found in exchange lookup", evt.Sender)
	}

	s0, err := ensureSymbol(ctx, db, evt.OfferAsset)
	if err != nil {
		return err
	}
	s1, err := ensureSymbol(ctx, db, evt.AskAsset)
	if err != nil {
		return err
	}

	i, err := ensureInstrument(ctx, db, &exchange, s0, s1)
	if err != nil {
		return err
	}

	price := calculatePrice(evt, i, s0, s1)
	pd := model.PriceData{
		InstrumentID: i.ID,
		Time:         evt.Timestamp,
		Price:        price,
	}

	// TODO: How do you handle multiple sawps for same pair in one TX?
	// The first one would insert, the second one would fail.
	_, err = db.NewInsert().Model(&pd).Exec(ctx)
	return err
}

func loadExchangeLkUp(ctx context.Context, db *persistence.Database) (map[string]model.Exchange, error) {
	var exchanges []model.ExchangeLkup
	if err := db.NewSelect().Model(&exchanges).Relation("Exchange").Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to fetch exchanges: %w", err)
	}

	lkup := make(map[string]model.Exchange)
	for _, v := range exchanges {
		lkup[v.Address] = *v.Exchange
	}

	return lkup, nil
}

func ensureSymbol(ctx context.Context, db bun.Tx, s string) (*model.Symbol, error) {
	var symbol model.Symbol
	if err := db.NewSelect().Model(&symbol).Where("id = ?", s).Scan(ctx); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to fetch symbol: %w", err)
		}
		sy, err := types.SymbolFromTokenDenom(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse symbol: %w", err)
		}
		symbol = model.Symbol{ID: s, DisplayName: sy.String()}
		if _, err := db.NewInsert().Model(&symbol).Exec(ctx); err != nil {
			return nil, err
		}
	}

	return &symbol, nil
}

func ensureInstrument(ctx context.Context, db bun.Tx, exchange *model.Exchange, s0, s1 *model.Symbol) (*model.Instrument, error) {
	if s0.ID == s1.ID {
		return nil, fmt.Errorf("symbols are identical")
	}

	// Ensure s0_id < s1_id to prevent double identity for same instrument
	sids := []string{s0.ID, s1.ID}
	slices.Sort(sids)

	var instrument model.Instrument
	stmt := db.NewSelect().Model(&instrument).
		Where("? = ?", bun.Ident("exchange_id"), exchange.ID).
		Where("? = ?", bun.Ident("symbol0_id"), sids[0]).
		Where("? = ?", bun.Ident("symbol1_id"), sids[1])

	if err := stmt.Scan(ctx); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to fetch instrument: %w", err)
		}

		ms0, err := types.SymbolFromTokenDenom(sids[0])
		if err != nil {
			return nil, err
		}
		ms1, err := types.SymbolFromTokenDenom(sids[1])
		if err != nil {
			return nil, err
		}
		iv := types.NewInstrument(ms0, ms1, exchange.Name)
		instrument = model.Instrument{
			Name:        iv.FullName(),
			ShortName:   iv.Name(),
			DisplayName: iv.Name(),
			Description: fmt.Sprintf("%s: [%s,%s]", exchange.Name, sids[0], sids[1]),
			ExchangeID:  exchange.ID,
			Symbol0ID:   sids[0],
			Symbol1ID:   sids[1],
		}
		if _, err := db.NewInsert().Model(&instrument).Exec(ctx); err != nil {
			return nil, err
		}
	}

	return &instrument, nil
}

func calculatePrice(evt types.SwapEvent, i *model.Instrument, s0, s1 *model.Symbol) float64 {
	a := evt.OfferAmount.Uint64()
	b := evt.ReturnAmount.Uint64()

	// this is to know how much for 1b
	price := float64(a) / float64(b)
	if i.Symbol0ID != s0.ID {
		price = 1 / price
	}

	return price
}
