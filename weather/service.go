package weather

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/drewthor/openweather/pkg/openweather"
)

func NewService(weatherClient openweather.Client) Service {
	return Service{weatherClient: weatherClient}
}

type Service struct {
	weatherClient openweather.Client
}

func (s Service) GetCurrentWeather(ctx context.Context, logger *slog.Logger, unitType UnitType, latitude float64, longitude float64) (WeatherData, error) {
	unit, err := getOpenWeatherUnit(unitType)
	if err != nil {
		return WeatherData{}, fmt.Errorf("invalid open weather unit: %w", err)
	}
	currentWeatherData, err := s.weatherClient.CurrentWeatherAtLocation(ctx, unit, latitude, longitude)
	if err != nil {
		return WeatherData{}, fmt.Errorf("failed to get current weather from client: %w", err)
	}

	weatherData := WeatherData{
		Temperature:     currentWeatherData.Main.Temp,
		UnitType:        unitType,
		Conditions:      weatherDescriptionsToConditions(currentWeatherData.WeatherConditions),
		TemperatureType: temperatureToType(currentWeatherData.Main.Temp, unitType),
	}

	return weatherData, nil
}

func temperatureToType(temperature float64, unitType UnitType) TemperatureType {
	tempFahrenheit := temperatureToFahrenheit(temperature, unitType)
	switch {
	case tempFahrenheit <= 40:
		return TemperatureTypeCold
	case tempFahrenheit >= 70:
		return TemperatureTypeHot
	default:
		return TemperatureTypeModerate
	}
}

func temperatureToFahrenheit(temperature float64, unitType UnitType) float64 {
	switch unitType {
	case UnitTypeMetric:
		return celsiusToFahrenheit(temperature)
	case UnitTypeScientific:
		return kelvinToFahrenheit(temperature)
	case UnitTypeImperial:
	default:
		return temperature
	}
	return temperature
}

func celsiusToFahrenheit(temperature float64) float64 {
	return (temperature * 1.8) + 32
}

func kelvinToFahrenheit(temperature float64) float64 {
	return celsiusToFahrenheit(temperature - 273.15)
}

func weatherDescriptionsToConditions(weatherConditions []openweather.WeatherCondition) []string {
	conditions := make([]string, 0, len(weatherConditions))
	for _, weatherCondition := range weatherConditions {
		conditions = append(conditions, strings.ToLower(weatherCondition.Description))
	}
	return conditions
}

func getOpenWeatherUnit(unitType UnitType) (openweather.Unit, error) {
	switch unitType {
	case UnitTypeScientific:
		return openweather.UnitStandard, nil
	case UnitTypeImperial:
		return openweather.UnitImperial, nil
	case UnitTypeMetric:
		return openweather.UnitMetric, nil
	default:
		return "", fmt.Errorf("cannot convert unit type to open weather unit")
	}
}
