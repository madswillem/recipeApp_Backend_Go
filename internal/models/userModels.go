package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/madswillem/recipeApp_Backend_Go/internal/error_handler"
	"github.com/madswillem/recipeApp_Backend_Go/internal/tools"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserModel struct {
	ID           string              `database:"id"`
	LastLogin    time.Time           `database:"last_login"`
	Cookie       string              `database:"cookie"`
	IP           string              `database:"ip"`
	RecipeGroups []RecipeGroupSchema `database:"groups"`
	Settings     UserSettings        `database:"settings"`
}
type UserSettings struct {
	Allergies []*IngredientDB `database:"allergies"`
	Diet      DietSchema      `gorm:"polymorphic:Owner"`
}

func (user *UserModel) GetByCookie(db *gorm.DB) *error_handler.APIError {
	err := db.Preload(clause.Associations).First(&user, "Cookie = ?", user.Cookie).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return error_handler.New("user not found", http.StatusNotFound, err)
		} else {
			return error_handler.New("database error", http.StatusInternalServerError, err)
		}
	}

	return nil
}

func (user *UserModel) CheckIfExistsByCookie(db *sqlx.DB) bool {
	found := false
	err := db.Get(found, "SELECT EXISTS(SELECT * FROM user WHERE cookie = $1) AS found;", user.Cookie)
	if err != nil {
		error_handler.New("database error", http.StatusInternalServerError, err)
	}
	return found
}

func (user *UserModel) Create(db *sqlx.DB, ip string) *error_handler.APIError {
	user.LastLogin = time.Now()
	user.IP = ip
	for {
		user.Cookie = tools.RandomString(20)
		if !user.CheckIfExistsByCookie(db) {
			break
		}
	}

	query := `INSERT INTO "user" (cookie, ip) VALUES (:cookie, :ip) RETURNING id`
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return error_handler.New("Query error: "+err.Error(), http.StatusInternalServerError, err)
	}
	err = stmt.Get(&user.ID, user)
	stmt.Close()
	if err != nil {
		return error_handler.New("Error inserting user: "+err.Error(), http.StatusInternalServerError, err)
	}

	return nil
}

func (user *UserModel) AddGroup(db *sqlx.DB, r *RecipeSchema) *error_handler.APIError {
	apiErr := r.GetRecipeByID(db)
	if apiErr != nil {
		return apiErr
	}

	rp := RecipeGroupSchema{}
	rp.Create(r)
	user.RecipeGroups = append(user.RecipeGroups, rp)

	v, err := json.Marshal(user.RecipeGroups)
	if err != nil {
		return error_handler.New("Failed to marshal recipe groups", http.StatusInternalServerError, err)
	}

	db.MustExec(`UPDATE "user" SET groups = $1 WHERE id = $2`, v, user.ID)
	return nil
}
func (user *UserModel) AddToGroup(db *sqlx.DB, r *RecipeSchema) *error_handler.APIError {
	var groups []byte
	err := db.Get(&groups, `SELECT "groups" FROM "user" WHERE id=$1`, user.ID)
	if err != nil {
		return error_handler.New("Error fetching recipe_groups", http.StatusInternalServerError, err)
	}
	json.Unmarshal(groups, &user.RecipeGroups)

	if len(user.RecipeGroups) < 1 {
		return user.AddGroup(db, r)
	}

	group_ranking := make([]struct {
		Group *RecipeGroupSchema
		Sim   float64
	}, len(user.RecipeGroups))

	for i := range user.RecipeGroups {
		group_ranking[i].Group = &user.RecipeGroups[i]
		group_ranking[i].Sim = user.RecipeGroups[i].Compare(r)
	}

	group_addble := make([]struct {
		Group *RecipeGroupSchema
		Sim   float64
	}, 0)
	fmt.Printf("%+v\n", group_ranking)
	for i := range group_ranking {
		if group_ranking[i].Sim >= .9 {
			group_addble = append(group_addble, group_ranking[i])
		}
	}
	fmt.Printf("%+v\n", group_ranking)

	if len(group_addble) > 1 {
		for i := 1; i < len(group_addble); i++ {
			group_addble[0].Group.Merge(group_addble[i].Group)
		}
	}

	group_addble[0].Group.Add(r)

	return nil
}

// Using db to extend an existing db like a search to show recipes similar to your intrests
func (user *UserModel) GetRecomendation(db *gorm.DB) (*error_handler.APIError, []RecipeSchema) {
	return nil, nil
}
