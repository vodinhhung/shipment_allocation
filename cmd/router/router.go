package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"shipment_allocation/internal/api"
	"shipment_allocation/internal/common"
	"shipment_allocation/internal/model"
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
	var data model.RequestData
	fmt.Println(r.Body)

	// Decode JSON body into struct
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println(&data)

	respData, err := router.allocation.ShipmentAllocation(&data)
	if err != nil {
		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// respond with allocation result
	common.WriteJSON(w, http.StatusOK, respData)
}

func NewRouter(allocation *api.Allocation) *Router {
	return &Router{
		allocation: allocation,
	}
}
