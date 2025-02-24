package service

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	env "github.com/caarlos0/env/v11"
	"github.com/palomachain/paloma-cdp/internal/pkg/liblog"
	"github.com/palomachain/paloma-cdp/internal/pkg/persistence"
)

type PersistenceRunner[T any] func(context.Context, Version, *persistence.Database, *T) error

type Version struct {
	Main   string
	Date   string
	Suffix string
}

func DefaultVersion() Version {
	return Version{
		Main:   "0.0.0",
		Date:   time.Now().Format(time.RFC3339),
		Suffix: "dev",
	}
}

func (v *Version) String() string {
	return fmt.Sprintf("%s-%s (%s)", v.Main, v.Suffix, v.Date)
}

type Shell[T any] struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	db        *persistence.Database
	version   Version
}

func New[T any]() *Shell[T] {
	liblog.Configure()
	ctx, fn := context.WithCancel(context.Background())
	return &Shell[T]{
		ctx:       ctx,
		ctxCancel: fn,
	}
}

func (s *Shell[T]) WithName(name string) *Shell[T] {
	s.ctx = liblog.HydrateServiceName(s.ctx, name)
	return s
}

func (s *Shell[T]) WithVersion(v Version) *Shell[T] {
	s.version = v
	return s
}

func (s *Shell[T]) WithDatabase() *Shell[T] {
	dbCfg, err := parseConfig[persistence.Configuration]()
	if err != nil {
		liblog.WithError(s.ctx, err, "failed to parse persistence config")
		panic(err)
	}

	db, err := persistence.New(s.ctx, dbCfg)
	if err != nil {
		liblog.WithError(s.ctx, err, "failed to connect to database")
		panic(err)
	}

	s.db = db
	return s
}

// RunWithPersistence runs the provided PersistenceRunner function with proper
// setup and teardown of the service, including handling OS signals for graceful
// shutdown and closing the database connection.
func (s *Shell[T]) RunWithPersistence(fn PersistenceRunner[T]) error {
	defer func() {
		slog.Default().InfoContext(s.ctx, "Service shutting down")
		if s.db != nil {
			if err := s.db.Close(); err != nil {
				liblog.WithError(s.ctx, err, "Failed to close database connection")
			}
		}
	}()

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	svcCfg, err := parseConfig[T]()
	if err != nil {
		return fmt.Errorf("failed to parse service config: %w", err)
	}

	go func() {
		errCh <- fn(s.ctx, s.version, s.db, svcCfg)
	}()

	select {
	case <-sigCh:
		slog.Default().InfoContext(s.ctx, "Shutdown signal received")
		s.ctxCancel()
		select {
		case err := <-errCh:
			if err != nil {
				liblog.WithError(s.ctx, err, "Failed to handle graceful shutdown")
			}
			return err
		case <-time.After(10 * time.Second):
			return nil
		}
	case err := <-errCh:
		if err != nil {
			liblog.WithError(s.ctx, err, "Service error encountered.")
		}
		return err
	}
}

func parseConfig[T any]() (*T, error) {
	r := new(T)
	return r, env.Parse(r)
}
