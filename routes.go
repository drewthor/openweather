package main

import (
	"net/http"

	"github.com/drewthor/openweather/weather"
)

func addRoutes(mux *http.ServeMux, weatherHandler weather.Handler) {
	mux.HandleFunc("GET /weather", weatherHandler.GetWeather)
}
