package tools

import (
	"encoding/json"
	"net/http"
	"time"
)

type IPStruct struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
type WeatherData struct {
	Latitude             float64        `json:"latitude"`
	Longitude            float64        `json:"longitude"`
	GenerationTimeMs     float64        `json:"generationtime_ms"`
	UTCOffsetSeconds     int            `json:"utc_offset_seconds"`
	Timezone             string         `json:"timezone"`
	TimezoneAbbreviation string         `json:"timezone_abbreviation"`
	Elevation            float64        `json:"elevation"`
	CurrentWeather       CurrentWeather `json:"current_weather"`
}
type CurrentWeather struct {
	Temperature   float64 `json:"temperature"`
	WindSpeed     float64 `json:"windspeed"`
	WindDirection float64 `json:"winddirection"`
	WeatherCode   int     `json:"weathercode"`
	Time          string  `json:"time"`
}
type CurrentData struct {
	Day    string
	Season string
	Temp   string
}

func getTemp() (float64, error) {
	url := "https://api.open-meteo.com/v1/forecast?latitude=53.5544&longitude=9.9946&current_weather=true"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "recipeapp")

	res, _ := http.DefaultClient.Do(req)

	var data WeatherData
	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	return data.CurrentWeather.Temperature, err
}

func GetCurrentData() (CurrentData, error) {
	var res CurrentData

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

	temp, err := getTemp()
	if err != nil {
		return res, err
	}

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

	return res, err
}
