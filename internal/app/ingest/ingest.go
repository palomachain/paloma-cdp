package ingest

import (
	"context"
	"fmt"
	"time"

	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cometbft/cometbft/types"
)

type Configuration struct {
	RpcAddress string `env:"PALOMA_RPC_ADDRESS,notEmpty"`
}

type Ingester struct {
	db  *persistence.Database
	cfg *Configuration
}

func NewIngester(db *persistence.Database, cfg *Configuration) *Ingester {
	return &Ingester{
		db:  db,
		cfg: cfg,
	}
}

func (i *Ingester) Run() error {
	client, err := rpchttp.New("https://rpc.palomachain.com:443")
	if err != nil {
		return err
	}
	err = client.Start()
	if err != nil {
		return err
	}
	defer client.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := "tm.event = 'Tx'"
	// query := "tm.event = 'Tx' AND tx.height = 3"
	// query := "tm.event = 'NewBlock'"
	txs, err := client.Subscribe(ctx, "test-client", query)
	if err != nil {
		return err
	}

	go func() {
		for e := range txs {
			data := e.Data.(types.EventDataTx)
			fmt.Println("got ", data)
		}
	}()

	fmt.Println("waiting for events")
	client.Wait()

	return nil
}
