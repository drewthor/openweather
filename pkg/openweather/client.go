package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strconv"
)

var (
	BaseURL = url.URL{Scheme: "https", Host: "api.openweathermap.org"}
)

type Unit string

const (
	UnitStandard Unit = "standard"
	UnitImperial Unit = "imperial"
	UnitMetric   Unit = "metric"
)

var (
	AllUnits = []Unit{UnitStandard, UnitImperial, UnitMetric}
)

func validUnit(u Unit) bool {
	return slices.Contains(AllUnits, u)
}

type Client struct {
	APIKey string
}

func NewClient(apiKey string) Client {
	return Client{APIKey: apiKey}
}

func (c Client) CurrentWeatherAtLocation(ctx context.Context, unit Unit, latitude float64, longitude float64) (CurrentWeatherData, error) {
	currentWeatherURLValues := url.Values{}
	currentWeatherURLValues.Add("lat", strconv.FormatFloat(latitude, 'g', -1, 64))
	currentWeatherURLValues.Add("lon", strconv.FormatFloat(longitude, 'g', -1, 64))
	currentWeatherURLValues.Add("appid", c.APIKey)
	if unit != "" {
		if !validUnit(unit) {
			return CurrentWeatherData{}, fmt.Errorf("invalid unit: %s", unit)
		}
		currentWeatherURLValues.Add("units", string(unit))
	}

	currentWeatherURL := BaseURL
	currentWeatherURL.Path = "/data/2.5/weather"
	currentWeatherURL.RawQuery = currentWeatherURLValues.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, currentWeatherURL.String(), nil)
	if err != nil {
		return CurrentWeatherData{}, fmt.Errorf("failed to create request to get current weather at location: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CurrentWeatherData{}, fmt.Errorf("failed to make request to get current weather at location: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentWeatherData{}, fmt.Errorf("failed to read response body getting current weather at location: %w", err)
	}

	currentWeather := CurrentWeatherData{}
	if err := json.Unmarshal(respBody, &currentWeather); err != nil {
		return CurrentWeatherData{}, fmt.Errorf("failed to unmarshal response body json getting current weather at location: %w", err)
	}

	return currentWeather, nil
}
