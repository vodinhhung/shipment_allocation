package common

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

func GetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ParseDurationOrDefault(v string, d time.Duration) time.Duration {
	if parsed, err := time.ParseDuration(v); err == nil {
		return parsed
	}
	return d
}

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}
