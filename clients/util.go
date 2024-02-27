package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anthonymq/go-stack-demo/logger"
	"go.uber.org/zap"
)

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

func unmarshal[T any](r *http.Response) (T, error) {
	var v T
	bodyBytes, err := io.ReadAll(r.Body)
	logger.Get().Debug("Response body", zap.String("body", string(bodyBytes)))
	if err != nil {
		logger.Get().Error("Error reading response body", zap.Error(err))
		return v, err
	}
	err = json.Unmarshal(bodyBytes, &v)
	if err != nil {
		logger.Get().Error("Unmarshall error", zap.Error(err))
		return v, err
	}

	return v, nil
}
