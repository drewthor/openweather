package openweather

type CurrentWeatherData struct {
	Coordinate struct {
		Longitude float64 `json:"lon"`
		Lattitude float64 `json:"lat"`
	} `json:"coord"`
	WeatherConditions []WeatherCondition `json:"weather"`
	Base              string             `json:"base"`
	Main              struct {
		Temp      float64 `json:"temp"` // default temp is in kelvin
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed  float64 `json:"speed"`
		Degree int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	DatePastEpoch int `json:"dt"`
	Sys           struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Code     int    `json:"cod"`
}

type WeatherCondition struct {
	ID          int    `json:"id"`
	Name        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
