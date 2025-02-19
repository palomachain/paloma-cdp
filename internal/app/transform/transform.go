package transform

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/uptrace/bun"
)

const (
	cPollingLimit int = 20
)

type Configuration struct {
	PollingInterval time.Duration `env:"CDP_TRANSFORMER_POLLING_INTERVAL,notEmpty" envDefault:"1s"`
}

func Run(
	ctx context.Context,
	db *persistence.Database,
	cfg *Configuration,
) error {
	slog.Default().InfoContext(ctx, "Service running.")

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
		exchangeLkup, err := loadExchangeLkUp(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to load exchange lookup: %w", err)
		}

		err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
			for _, v := range txs {
				if err := consume(ctx, db, exchangeLkup, v); err != nil {
					return err
				}
			}

			// TODO: Insert & Delete
			return nil
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

func consume(ctx context.Context,
	db *persistence.Database,
	exchangeLkup map[string]model.Exchange,
	tx model.RawTransaction,
) error {
	slog.Default().InfoContext(ctx, "Consuming transaction", "hash", tx.Hash)

	for _, v := range tx.Data.Result.Events {
		fmt.Printf("Event type: %v\n", v.Type)
		for _, j := range v.Attributes {
			fmt.Printf("-- Attribute: %v\n", j)
		}
	}

	return nil

	if !isSwapTx(ctx, tx) {
		slog.Default().DebugContext(ctx, "Skipping non-swap tx.", "hash", tx.Hash)
		return nil
	}

	return nil
}

func isSwapTx(ctx context.Context, tx model.RawTransaction) bool {
	// Offer ammount & offer asset are EXACTLY what was sent to the contract
	// Return amount & asset are EXACTLY what the RECEIVER gets back
	//
	// Sender should be paloma17nm703yu6vy6jpwn686e5ucal7n4cw8fc6da9ee0ctcwmr9vc9nsr4evrh for
	// Bonding curve or the other one
	for range []string{
		"wasm._contract_address",
		"wasm.action",
		"wasm.sender",
		"wasm.receiver",
		"wasm.ask_asset",
		"wasm.offer_asset",
		"wasm.offer_amount",
		"wasm.return_amount",
	} {
		// i, ok := tx..Data[v]
		// if !ok || len(i) < 1 {
		// 	return false
		// }
	}

	return true
}

func loadExchangeLkUp(ctx context.Context, db *persistence.Database) (lkup map[string]model.Exchange, err error) {
	var exchanges []model.ExchangeLkup
	if err = db.NewSelect().Model(&exchanges).Relation("Exchange").Scan(ctx); err != nil {
		return nil, fmt.Errorf("failed to fetch exchanges: %w", err)
	}

	// for _, v := range exchanges {
	// 	lkup[v.Address] = *v.Exchange
	// }

	return lkup, nil
}
