package container

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainer struct {
	postgres.PostgresContainer
}

func StartPostgres(ctx context.Context) (*PostgresContainer, error) {
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

