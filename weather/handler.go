package weather

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func NewHandler(logger *slog.Logger, service Service) Handler {
	return Handler{logger: logger, service: service}
}

type Handler struct {
	logger  *slog.Logger
	service Service
}

func (h Handler) GetWeather(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		h.logger.ErrorContext(ctx, "failed to parse http request form for get weather", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	latitudeStr := r.FormValue("latitude")
	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		b, _ := json.Marshal(jsonError{Error: "invalid latitude; must be number"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	longitudeStr := r.FormValue("longitude")
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		b, _ := json.Marshal(jsonError{Error: "invalid longitude; must be number"})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	unitType, err := parseUnit(r.FormValue("unit"))
	if err != nil {
		b, _ := json.Marshal(jsonError{Error: fmt.Sprintf("invalid unit type; %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	weatherData, err := h.service.GetCurrentWeather(ctx, h.logger, unitType, latitude, longitude)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to get current weather", slog.Any("error", err))
		b, _ := json.Marshal(jsonError{Error: fmt.Sprintf("invalid unit type; %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	b, err := json.Marshal(weatherData)
	if err != nil {
		h.logger.ErrorContext(ctx, "failed to marshal json response for current weather", slog.Any("error", err))
		b, _ := json.Marshal(jsonError{Error: fmt.Sprintf("invalid unit type; %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	if bytesWritten, err := w.Write(b); err != nil {
		h.logger.ErrorContext(ctx, fmt.Sprintf("failed to write weather data response: wrote %d bytes", bytesWritten), slog.Any("error", err))
		return
	}
}

func parseUnit(unitStr string) (UnitType, error) {
	switch unitStr {
	case "":
	case string(UnitTypeScientific):
		return UnitTypeScientific, nil
	case string(UnitTypeMetric):
		return UnitTypeMetric, nil
	case string(UnitTypeImperial):
		return UnitTypeImperial, nil
	default:
		return "", fmt.Errorf(fmt.Sprintf("unit type must be one of [%s]", AllUnitTypes))
	}
	return "", fmt.Errorf(fmt.Sprintf("unit type must be one of %v", AllUnitTypes))
}

type jsonError struct {
	Error string `json:"error"`
}
