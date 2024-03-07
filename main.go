package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/drewthor/openweather/pkg/openweather"
	"github.com/drewthor/openweather/weather"
)

type config struct {
	Port              int
	OpenWeatherAPIKey string
}

func main() {
	var cfg config

	cfg.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	cfg.OpenWeatherAPIKey = os.Getenv("OPEN_WEATHER_API_KEY")

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()

	/*if cfg.enableTracing {
		// Configure a new exporter using environment variables for sending data to Honeycomb over gRPC.
		exporter, err := otlptracegrpc.New(ctx)
		if err != nil {
			log.Fatalf("failed to initialize exporter: %v", err)
		}

		// Create a new tracer provider with a batch span processor and the otlp exporter.
		tp := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(exporter),
		)

		// Handle shutdown errors in a sensible manner where possible
		defer func() { _ = tp.Shutdown(ctx) }()

		// Set the Tracer Provider global
		otel.SetTracerProvider(tp)

		// Register the trace context and baggage propagators so data is propagated across services/processes.
		otel.SetTextMapPropagator(
			propagation.NewCompositeTextMapPropagator(
				propagation.TraceContext{},
				propagation.Baggage{},
			),
		)
	}*/

	weatherClient := openweather.NewClient(cfg.OpenWeatherAPIKey)
	weatherHandler := weather.NewHandler(logger, weather.NewService(weatherClient))

	mux := http.NewServeMux()
	addRoutes(mux, weatherHandler)

	// srvHandler := middleware.Recoverer(handler.Routes())

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		TLSConfig:    nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
		ErrorLog:     log.Default(), // use default logger set by slog.SetDefault()
	}

	logger.Info(fmt.Sprintf("starting server on port %d", cfg.Port))
	err := srv.ListenAndServe()
	logger.ErrorContext(ctx, "server shutdown", slog.Any("error", err.Error()))
}
