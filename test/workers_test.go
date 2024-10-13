package test

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/database"
	"github.com/madswillem/recipeApp_Backend_Go/internal/workers"
)

func TestCreateSelectedAndViewdLog(t *testing.T) {
	container, ctx := InitTestContainer(t)
	URL, err := container.ConnectionString(*ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	db := database.ConnectToDB(&sqlx.Conn{}, URL)

	err = workers.CreatSelectedAndViewLog(db)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(workers.GetLastLog(db))

}
