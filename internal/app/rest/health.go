package rest

import (
	"context"

	"github.com/swaggest/usecase"
)

func HealthInteractor() usecase.IOInteractor {
	u := usecase.NewInteractor(func(ctx context.Context, _, _ *struct{}) error {
		return nil
	})

	u.SetTitle("Health")
	u.SetDescription("The health probe endpoint can be used to determine if the service is healthy.")
	u.SetTags("Core")

	return u.IOInteractor
}
