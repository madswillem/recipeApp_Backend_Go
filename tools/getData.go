package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"rezeptapp.ml/goApp/models"
)

func getTemp(c *gin.Context) float64 {
	url := "https://api.open-meteo.com/v1/forecast?latitude=53.5544&longitude=9.9946&current_weather=true"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "recipeapp")

	res, _ := http.DefaultClient.Do(req)

	var data models.WeatherData
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println("Error:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	return data.CurrentWeather.Temperature
}

func GetCurrentData(c *gin.Context) models.CurrentData {
	var res models.CurrentData

	month := time.Now().Month()
	var season string
	switch month {
	case time.December, time.January, time.February:
		season = "Win"
	case time.March, time.April, time.May:
		season = "Spr"
	case time.June, time.July, time.August:
		season = "Sum"
	case time.September, time.October, time.November:
		season = "Aut"
	default:
		season = "non"
	}

	temp := getTemp(c)
	switch {
	case temp <= 0.0:
		res.Temp = "subzerodegree"
	case temp <= 10.0 && temp > 0.0:
		res.Temp = "zerodegree"
	case temp <= 20.0 && temp > 10.0:
		res.Temp = "tendegree"
	case temp <= 30.0 && temp > 20.0:
		res.Temp = "twentiedegree"
	case temp > 30.0:
		res.Temp = "thirtydegree"
	}

	res.Season = season
	res.Day = time.Now().Weekday().String()[:3]

	return res
}