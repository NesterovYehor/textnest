package container

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type PostgresContainerOptions struct {
	Database string
	Host     string
	UserName string
	Password string
}

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

	t.Cleanup(func() {
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	return &PostgresContainer{*postgresContainer}, nil
}
