package models

import (
	"time"

	"gorm.io/gorm"
)

type RatingStruct struct {
	ID		  	uint   `gorm:"primarykey"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	OwnerTitle	string `json:"owner_title"`
	OwnerID   	string
  	OwnerType 	string

	Overall float64 `json:"overall"`

	Mon float64 `json:"mon"`
	Tue float64 `json:"tue"`
	Wed float64 `json:"wed"`
	Thu float64 `json:"thu"`
	Fri float64 `json:"fri"`
	Sat float64 `json:"sat"`
	Sun float64 `json:"sun"`

	Win float64 `json:"win"`
	Spr float64 `json:"spr"`
	Sum float64 `json:"sum"`
	Aut float64 `json:"aut"`

	Thirtydegree  float64 `json:"thirtydegree"`
	Twentiedegree float64 `json:"twentiedegree"`
	Tendegree     float64 `json:"tendegree"`
	Zerodegree    float64 `json:"zerodegree"`
	Subzerodegree float64 `json:"subzerodegree"`
}

// constructor function
func NewRatingStruct(title string) *RatingStruct {
	return &RatingStruct{
		OwnerTitle:	  title,
		Overall:       1000,
		Mon:           1000,
		Tue:           1000,
		Wed:           1000,
		Thu:           1000,
		Fri:           1000,
		Sat:           1000,
		Sun:           1000,
		Win:           1000,
		Spr:           1000,
		Sum:           1000,
		Aut:           1000,
		Thirtydegree:  1000,
		Twentiedegree: 1000,
		Tendegree:     1000,
		Zerodegree:    1000,
		Subzerodegree: 1000,
	}
}
