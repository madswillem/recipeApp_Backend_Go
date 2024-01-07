package models

type RatingStruct struct {
	ID		  	uint   `gorm:"primarykey"`
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

func (rating *RatingStruct) DefaultRatingStruct(title string){
	rating.OwnerTitle = title

	rating.Overall = 1000.0
	rating.Mon = 1000.0
	rating.Tue = 1000.0
	rating.Wed = 1000.0
	rating.Thu = 1000.0
	rating.Fri = 1000.0
	rating.Sat = 1000.0
	rating.Sun = 1000.0
	rating.Win = 1000.0
	rating.Spr = 1000.0
	rating.Sum = 1000.0
	rating.Aut = 1000.0
	rating.Thirtydegree = 1000.0
	rating.Twentiedegree = 1000.0
	rating.Tendegree = 1000.0
	rating.Zerodegree = 1000.0
	rating.Subzerodegree = 1000.0
}
