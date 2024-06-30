package models

import (
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
)

type IngredientDB struct {
    ID           string    `db:"id" json:"id"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    Name         string    `db:"name" json:"name"`
    StandardUnit string   `db:"standard_unit" json:"standard_unit,omitempty"`
    NdbNumber    int64    `db:"ndb_number" json:"ndb_number,omitempty"`
    Category     string   `db:"category" json:"category,omitempty"`
    FdicID       int64    `db:"fdic_id" json:"fdic_id,omitempty"`
	NutritionalValue NutritionalValue `json:"nv"`
	Rating	 RatingStruct `json:"rating"`
}

func (ingredient *IngredientDB) Create(db *sqlx.Tx) (string, *error_handler.APIError) {
    
    return "", nil
}

func GetIngIDByName(tx *sqlx.Tx, name string) (string, *error_handler.APIError) {
    var id string
    err := tx.QueryRow("SELECT id FROM ingredient WHERE name = $1", name).Scan(id)
    if err != nil {
        if err.Error() == "sql: no rows in result set"{
            ing := IngredientDB{Name: name}
            return ing.Create(tx)
        }
        return "", error_handler.New("database error getting Ing: " + err.Error(), http.StatusInternalServerError, err)
    }

    return id, nil
}