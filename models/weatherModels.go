package models

type WeatherData struct {
	Latitude 				float64 			`json:"latitude"`
	Longitude 				float64 			`json:"longitude"`
	GenerationTimeMs 		float64 			`json:"generationtime_ms"`
	UTCOffsetSeconds 		int 				`json:"utc_offset_seconds"`
	Timezone 				string 				`json:"timezone"`
	TimezoneAbbreviation	string 				`json:"timezone_abbreviation"`
	Elevation 				float64 			`json:"elevation"`
	CurrentWeather 			CurrentWeather 		`json:"current_weather"`
}
type CurrentWeather struct {
	Temperature 			float64 			`json:"temperature"`
	WindSpeed 				float64 			`json:"windspeed"`
	WindDirection 			float64 			`json:"winddirection"`
	WeatherCode 			int 				`json:"weathercode"`
	Time 					string 				`json:"time"`
}