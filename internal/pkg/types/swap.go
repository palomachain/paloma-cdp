package types

import (
	"context"
	"fmt"
	"math/big"
	"time"

	v1 "github.com/cometbft/cometbft/api/cometbft/abci/v1"
)

type SwapEvent struct {
	ContractAddress string
	Action          string
	Sender          string
	Receiver        string
	AskAsset        string
	OfferAsset      string
	OfferAmount     *big.Int
	ReturnAmount    *big.Int
	Timestamp       time.Time // Used only to carry forward the TX received timestamp
}

func (e *SwapEvent) String() string {
	return fmt.Sprintf("%s %s%s to %s", e.Sender, e.OfferAmount.String(), e.OfferAsset, e.Receiver)
}

type CoinReceivedEvent struct {
	Amount   string
	Receiver string
}

// Offer ammount & offer asset are EXACTLY what was sent to the contract
// Return amount & asset are EXACTLY what the RECEIVER gets back
//
// Sender should be paloma17nm703yu6vy6jpwn686e5ucal7n4cw8fc6da9ee0ctcwmr9vc9nsr4evrh for
// Bonding curve or the other one
var swapMap = map[string]func(e *SwapEvent, i string){
	"_contract_address": func(e *SwapEvent, i string) { e.ContractAddress = i },
	"action":            func(e *SwapEvent, i string) { e.Action = i },
	"sender":            func(e *SwapEvent, i string) { e.Sender = i },
	"receiver":          func(e *SwapEvent, i string) { e.Receiver = i },
	"ask_asset":         func(e *SwapEvent, i string) { e.AskAsset = i },
	"offer_asset":       func(e *SwapEvent, i string) { e.OfferAsset = i },
	"offer_amount":      func(e *SwapEvent, i string) { e.OfferAmount, _ = big.NewInt(0).SetString(i, 10) },
	"return_amount":     func(e *SwapEvent, i string) { e.ReturnAmount, _ = big.NewInt(0).SetString(i, 10) },
}

var crMap = map[string]func(e *CoinReceivedEvent, i string){
	"receiver": func(e *CoinReceivedEvent, i string) { e.Receiver = i },
	"amount":   func(e *CoinReceivedEvent, i string) { e.Amount = i },
}

func (s *SwapEvent) HasSwapAction() bool {
	return s.Action == "swap"
}

func (s *SwapEvent) Validate() error {
	if s.OfferAmount == nil || s.ReturnAmount == nil {
		return fmt.Errorf("OfferAmount or ReturnAmount is nil")
	}

	if s.OfferAmount.Cmp(big.NewInt(0)) < 1 || s.ReturnAmount.Cmp(big.NewInt(0)) < 1 {
		return fmt.Errorf("OfferAmount or ReturnAmount is less than 1")
	}

	if s.ContractAddress == "" || s.Sender == "" || s.Receiver == "" || s.AskAsset == "" || s.OfferAsset == "" {
		return fmt.Errorf("ContractAddress, Sender, Receiver, AskAsset, or OfferAsset is empty")
	}

	return nil
}

func TryParseSwapEvents(ctx context.Context, events []v1.Event) ([]SwapEvent, error) {
	e := make([]SwapEvent, 0, len(events))
	cr := make([]CoinReceivedEvent, 0, len(events))

	for _, v := range events {
		if v.Type == "coin_received" {
			evt, err := tryParseCoinReceivedEvent(ctx, v.Attributes)
			if err != nil {
				return nil, fmt.Errorf("failed to parse coin received event: %w", err)
			}
			cr = append(cr, evt)
			continue
		}

		if v.Type != "wasm" {
			continue
		}

		evt := tryParseSwapEvent(ctx, v.Attributes)
		if !evt.HasSwapAction() {
			continue
		}

		if err := evt.Validate(); err != nil {
			return nil, fmt.Errorf("failed to validate swap event: %w", err)
		}

		if err := tryMatchSwapToCoinEvents(ctx, e, cr); err != nil {
			return nil, fmt.Errorf("failed to match swap to coin received events: %w", err)
		}

		e = append(e, evt)
	}

	return e, nil
}

func tryMatchSwapToCoinEvents(ctx context.Context, se []SwapEvent, cr []CoinReceivedEvent) error {
	pool := cr[:]

	for _, v := range se {
		for i, c := range pool {
			if c.Receiver == v.Sender {
				amount := fmt.Sprintf("%s%s", v.OfferAmount.String(), v.OfferAsset)
				if c.Amount == amount {
					pool = append(pool[:i], pool[i+1:]...)
					break
				}
			}
			return fmt.Errorf("failed to match swap event %s to any coin received event", v.String())
		}
	}

	return nil
}

func tryParseSwapEvent(ctx context.Context, a []v1.EventAttribute) SwapEvent {
	e := SwapEvent{}
	for _, v := range a {
		fn, ok := swapMap[v.Key]
		if !ok {
			continue
		}

		fn(&e, v.Value)
	}

	return e
}

func tryParseCoinReceivedEvent(ctx context.Context, a []v1.EventAttribute) (CoinReceivedEvent, error) {
	e := CoinReceivedEvent{}
	for _, v := range a {
		fn, ok := crMap[v.Key]
		if !ok {
			continue
		}

		fn(&e, v.Value)
	}

	if e.Receiver == "" || e.Amount == "" {
		return e, fmt.Errorf("Receiver or Amount is empty")
	}

	return e, nil
}
