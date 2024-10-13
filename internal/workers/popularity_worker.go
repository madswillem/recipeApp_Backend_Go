package workers

import (
	"github.com/jmoiron/sqlx"
)

type r_log struct {
	ID             string `db:"recipe_id"`
	Selects        int    `db:"selects"`
	Views          int    `db:"views"`
	View_change    int    `db:"view_change"`
	Selects_change int    `db:"selects_change"`
}

func CreatSelectedAndViewLog(db *sqlx.DB) error {
	r := []r_log{}
	err := db.Select(&r, "SELECT id as recipe_id, selects, views FROM recipes")
	if err != nil {
		print("hi")
		return err
	}
	r_l := []r_log{}
	r_l, err = GetLastLog(db)
	if err != nil {
		print("ho")
		return err
	}

	r_n := CreateDiff(r, r_l)
	_, err = db.NamedExec(`INSERT INTO recipe_selects_views_log (recipe_id, selects, views, view_change, selects_change)
		VALUES (:recipe_id, :selects, :views, :view_change, :selects_change)`, r_n)

	return err
}

func GetLastLog(db *sqlx.DB) ([]r_log, error) {
	r := []r_log{}
	err := db.Select(&r, `SELECT DISTINCT ON (recipe_id) recipe_id, selects, views, view_change, selects_change
			FROM recipe_selects_views_log
			ORDER BY recipe_id, day DESC;`)
	return r, err
}

func CreateDiff(r1, r2 []r_log) []r_log {
	var d []r_log
	for _, v1 := range r1 {
		found := false
		for _, v2 := range r2 {
			if v1 == v2 {
				d = append(d, r_log{ID: v1.ID, Selects: v1.Selects, Views: v1.Views, Selects_change: v1.Selects - v2.Selects, View_change: v1.Views - v2.Views})
				found = true
				break
			}
		}
		if !found {
			d = append(d, v1)
		}
	}

	return d
}
