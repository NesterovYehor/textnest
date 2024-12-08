package container

import (
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	postgres.PostgresContainer
}

func Start(ctx context.Context, t *testing.T) (*PostgresContainer, error) {
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test_db"),
		postgres.WithUsername("testcontainer"),
		postgres.WithPassword("testcontainer"),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		return nil, err
	}
	return &PostgresContainer{*postgresContainer}, nil
}

func PrintPenis()(){
    fmt.Println("PENIS")
}
