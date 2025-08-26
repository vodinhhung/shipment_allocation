package router

import (
	"net/http"
	"shipment_allocation/internal/common"
	"time"
)

func HandleRoot(w http.ResponseWriter, r *http.Request) {
	common.WriteJSON(w, http.StatusOK, map[string]string{"message": "hello from Go server"})
}

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	common.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}
