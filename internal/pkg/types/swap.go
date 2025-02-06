package types

import (
	"context"
	"fmt"
	"math/big"

	"github.com/palomachain/paloma-cdp/internal/pkg/model"
)

// Offer ammount & offer asset are EXACTLY what was sent to the contract
// Return amount & asset are EXACTLY what the RECEIVER gets back
//
// Sender should be paloma17nm703yu6vy6jpwn686e5ucal7n4cw8fc6da9ee0ctcwmr9vc9nsr4evrh for
// Bonding curve or the other one
var parseMap = map[string]func(e *SwapEvent, i string){
	"wasm._contract_address": func(e *SwapEvent, i string) { e.ContractAddress = i },
	"wasm.action":            func(e *SwapEvent, i string) { e.Action = i },
	"was.sender":             func(e *SwapEvent, i string) { e.Sender = i },
	"wasm.receiver":          func(e *SwapEvent, i string) { e.Receiver = i },
	"wasm.ask_asset":         func(e *SwapEvent, i string) { e.AskAsset = i },
	"wasm.offer_asset":       func(e *SwapEvent, i string) { e.OfferAsset = i },
	"wasm.offer_amount":      func(e *SwapEvent, i string) { e.OfferAmount, _ = big.NewInt(0).SetString(i, 10) },
	"wasm.return_amount":     func(e *SwapEvent, i string) { e.ReturnAmount, _ = big.NewInt(0).SetString(i, 10) },
}

type SwapEvent struct {
	ContractAddress string
	Action          string
	Sender          string
	Receiver        string
	AskAsset        string
	OfferAsset      string
	OfferAmount     *big.Int
	ReturnAmount    *big.Int
}

// TODO: What about if you have:
// 1 WASM event, not swap (but with _contract_address)
// 2 SWAP events
// or a wild mix
func tryParseSwapEvents(ctx context.Context, events map[string][]string) (*SwapEvent, error) {
	e := &SwapEvent{}
	for k, fn := range parseMap {
		i, ok := events[k]
		if !ok || len(i) < 1 {
			return nil, fmt.Errorf("missing key: %s", k)
		}
		fn(e, i[0])
	}
}

func isSwapTx(ctx context.Context, tx model.RawTransaction) bool {
	for _, v := range []string{
		"wasm._contract_address",
		"wasm.action",
		"wasm.sender",
		"wasm.receiver",
		"wasm.ask_asset",
		"wasm.offer_asset",
		"wasm.offer_amount",
		"wasm.return_amount",
	} {
		i, ok := tx.Data[v]
		if !ok || len(i) < 1 {
			return false
		}
	}

	return true
}
