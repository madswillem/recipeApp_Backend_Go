package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/server"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestGetRecipeByID(t *testing.T) {
	ctx := context.Background()

	// 1. Start the postgres container and run any migrations on it
	container, err := postgres.Run(
		ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithDatabase("test"),
		postgres.WithUsername("mads"),
		postgres.WithPassword("1234"),
		postgres.BasicWaitStrategies(),
		postgres.WithInitScripts("./testdata/innit-db.sql"),
		postgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create a snapshot of the database to restore later
	err = container.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	dbURL, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s := server.Server{NewDB: database.ConnectToDB(&sqlx.Conn{}, dbURL)}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Request, _ = http.NewRequest(http.MethodGet, "/get", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	s.GetAll(c)
}
