package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
	"github.com/palomachain/paloma-cdp/internal/pkg/service"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cometbft/cometbft/types"
)

type Configuration struct {
	RpcAddress string `env:"CDP_PALOMA_RPC_ADDRESS,notEmpty"`
	Query      string `env:"CDP_SUBSCRIPTION_QUERY,notEmpty" envDefault:"tm.event = 'Tx' AND message.action = '/cosmwasm.wasm.v1.MsgExecuteContract'"`
}

func Run(
	ctx context.Context,
	v service.Version,
	db *persistence.Database,
	cfg *Configuration,
) error {
	client, err := rpchttp.NewWithTimeout(cfg.RpcAddress, 5)
	if err != nil {
		return err
	}

	err = client.Start()
	if err != nil {
		return err
	}
	defer client.Stop()

	tk := time.NewTicker(1 * time.Second)

	txs, err := client.Subscribe(ctx, "", cfg.Query)
	if err != nil {
		return err
	}

	slog.Default().InfoContext(ctx, "Service running.", "query", cfg.Query, "version", v)
	for {
		select {
		case <-ctx.Done():
			tk.Stop()
			return client.UnsubscribeAll(ctx, "")
		case tx := <-txs:
			if err := handleTx(ctx, db, tx); err != nil {
				liblog.WithError(ctx, err, "Failed to handle tx", "events", tx.Events)
			}
		case <-tk.C:
			if !client.IsRunning() {
				slog.Default().WarnContext(ctx, "WSEvents not running. Trying to recover...")
				if err := client.Reset(); err != nil {
					liblog.WithError(ctx, err, "Failed to reset WS connection.")
				}
				if client.IsRunning() {
					slog.Default().InfoContext(ctx, "WSEvents recovered. Do we need to subscribe again?")
				}
			}
		}
	}
}

func serialize(data map[string][]string) ([]byte, error) {
	return json.Marshal(data)
}

func deserialize(data []byte) (result map[string][]string, err error) {
	return result, json.Unmarshal(data, &result)
}

func handleTx(ctx context.Context, db *persistence.Database, tx coretypes.ResultEvent) error {
	hash, err := tryGetValue(tx, "tx.hash")
	if err != nil {
		return err
	}

	data := tx.Data.(types.EventDataTx)

	m := &model.RawTransaction{
		Hash: hash,
		Data: data,
	}

	_, err = db.NewInsert().Model(m).Exec(ctx)
	if err != nil {
		return err
	}

	slog.Default().InfoContext(ctx, "Ingested transaction.", "hash", hash)

	return nil
}

func tryGetValue(tx coretypes.ResultEvent, key string) (string, error) {
	if v, ok := tx.Events[key]; ok {
		if len(v) < 1 {
			return "", fmt.Errorf("no value present for existing key %s", key)
		}
		return v[0], nil
	}
	return "", fmt.Errorf("key %s not found", key)
}

func _dbg_handleTx(tx coretypes.ResultEvent) {
	fmt.Printf("Hash: %v\n", tx.Events["tx.hash"][0])
	fmt.Println("====================")
	for k, v := range tx.Events {
		fmt.Printf("%s: %v\n", k, v)
	}
	fmt.Println("-----------------")
	data := tx.Data.(types.EventDataTx)
	for _, evt := range data.Result.Events {
		fmt.Printf("%v\n", evt)
	}
}
