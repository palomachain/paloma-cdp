package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/palomachain/paloma-cdp/internal/app/config"
	"github.com/palomachain/paloma-cdp/internal/pkg/model"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
)

// atm, approx 1GB per 20 symbols over 1m of 5 second data

const (
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numExchanges = 2
	numSymbols   = 100
)

var exchanges = []string{"PALOMADEX", "CURVEBOND"}

func randomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func randomPick[T any](slice []T) T {
	return slice[rand.Intn(len(slice))]
}

func main() {
	os.Setenv("CDP_PSQL_ADDRESS", "localhost:5432")
	os.Setenv("CDP_PSQL_USER", "cdp")
	os.Setenv("CDP_PSQL_PASSWORD", "trustno1")
	os.Setenv("CDP_PSQL_DATABASE", "cdp")
	os.Setenv("CDP_PSQL_TIMEOUT", "60s")
	os.Setenv("CDP_GQL_PORT", "8080")
	ctx := context.Background()

	cfg, err := config.Parse()
	if err != nil {
		slog.Default().ErrorContext(ctx, "failed to parse config: %v", err)
		panic(err)
	}

	db, err := persistence.New(ctx, &cfg.Persistence)
	if err != nil {
		slog.Default().ErrorContext(ctx, "failed to connect to database: %v", err)
		panic(err)
	}

	if err := db.Migrate(ctx); err != nil {
		slog.Default().ErrorContext(ctx, "failed to migrate database: %v", err)
		panic(err)
	}

	for _, t := range []any{
		(*model.PriceData)(nil),
		(*model.Symbol)(nil),
		(*model.Exchange)(nil),
	} {

		_, err := db.NewDelete().Model(t).Where("1=1").Exec(ctx)
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < numExchanges; i++ {
		m := model.Exchange{}
		m.Name = exchanges[i]
		if _, err := db.NewInsert().Model(&m).Exec(ctx); err != nil {
			panic(err)
		}
	}

	exchanges := []model.Exchange{}
	db.NewSelect().Model(&model.Exchange{}).Scan(ctx, &exchanges)

	for i := 0; i < numSymbols; i++ {
		m := model.Symbol{}
		m.Name = symbolName()
		m.Description = randomString(100)
		m.ExchangeID = randomPick(exchanges).ID
		if _, err := db.NewInsert().Model(&m).Exec(ctx); err != nil {
			panic(err)
		}
	}

	symbols := []model.Symbol{}
	db.NewSelect().Model(&model.Symbol{}).Scan(ctx, &symbols)

	to := time.Now().UTC()
	from := to.Add(-time.Hour * 24 * 30)

	for i, s := range symbols {
		fmt.Printf("[%d/%d]Inserting for %v\n", i, len(symbols), s.ID)
		idx := from

		prices := make([]model.PriceData, 0, 518400)
		for idx.Before(to) {
			m := model.PriceData{}
			m.SymbolID = s.ID
			m.Time = idx
			m.Price = rand.Float64()
			prices = append(prices, m)
			idx = idx.Add(time.Second * 5)
		}

		fmt.Printf("Len %d\n", len(prices))
		if _, err := db.NewInsert().Model(&prices).Exec(ctx); err != nil {
			panic(err)
		}
	}
}

var adjectives = []string{
	"swift", "brave", "lucky", "happy", "bright", "clever", "bold", "eager", "gentle", "jolly",
	"kind", "lively", "merry", "nimble", "proud", "quick", "silly", "tender", "witty", "zany",
	"fierce", "cheerful", "graceful", "shiny", "quirky", "bouncy", "daring", "curious", "dreamy", "sassy",
	"snappy", "breezy", "cuddly", "dandy", "feisty", "goofy", "jazzy", "peppy", "spunky", "whimsical",
	"snazzy", "zesty", "giddy", "perky", "waggish", "jaunty", "quirky", "sprightly", "vivid", "zealous",
}

var nouns = []string{
	"panda", "lion", "tiger", "eagle", "shark", "wolf", "otter", "beaver", "dolphin", "hawk",
	"rabbit", "penguin", "falcon", "turtle", "koala", "fox", "squirrel", "raccoon", "whale", "stingray",
	"octopus", "owl", "buffalo", "zebra", "cheetah", "giraffe", "lemur", "armadillo", "pigeon", "peacock",
	"toucan", "jaguar", "panther", "hippo", "kangaroo", "mongoose", "reindeer", "sloth", "tapir", "vulture",
	"yak", "wolverine", "porcupine", "meerkat", "boar", "chameleon", "gazelle", "lynx", "badger", "jackal",
}

func symbolName() string {
	address := fmt.Sprintf("paloma1%s", randomString(38))
	return fmt.Sprintf("tokenfactory/%s/%s%s", address, randomPick(adjectives), randomPick(nouns))
}
