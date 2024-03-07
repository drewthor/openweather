package weather

import "context"

type Client interface {
	CurrentWeatherAtLocation(ctx context.Context, lattitude float64, longitude float64) (WeatherData, error)
}

type WeatherData struct {
	Temperature     float64         `json:"temperature"`
	UnitType        UnitType        `json:"unit_type"`
	Conditions      []string        `json:"conditions"`
	TemperatureType TemperatureType `json:"temperature_type"`
}

type UnitType string

const (
	UnitTypeScientific UnitType = "scientific"
	UnitTypeImperial   UnitType = "imperial"
	UnitTypeMetric     UnitType = "metric"
)

var (
	AllUnitTypes = []UnitType{UnitTypeScientific, UnitTypeImperial, UnitTypeMetric}
)

type TemperatureType string

const (
	TemperatureTypeHot      TemperatureType = "hot"
	TemperatureTypeCold     TemperatureType = "cold"
	TemperatureTypeModerate TemperatureType = "moderate"
)
