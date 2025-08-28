package router

import (
	"net/http"
	"shipment_allocation/internal/api"
	"shipment_allocation/internal/common"
	"time"
)

type Router struct {
	allocation *api.Allocation
}

func (router *Router) HandleRoot(w http.ResponseWriter, r *http.Request) {
	common.WriteJSON(w, http.StatusOK, map[string]string{"message": "hello from Go server"})
}

func (router *Router) HandleHealth(w http.ResponseWriter, r *http.Request) {
	common.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "time": time.Now().UTC().Format(time.RFC3339)})
}

func (router *Router) HandleAllocateShipment(w http.ResponseWriter, r *http.Request) {
	err := router.allocation.ShipmentAllocation()
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := map[string]string{"status": "ok"}

	// respond with allocation result
	common.WriteJSON(w, http.StatusOK, resp)
}

func NewRouter(allocation *api.Allocation) *Router {
	return &Router{
		allocation: allocation,
	}
}
