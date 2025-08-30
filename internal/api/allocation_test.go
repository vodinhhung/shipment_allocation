package api

import (
	"shipment_allocation/internal/dependency"
	"shipment_allocation/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolveILP_BasicCase(t *testing.T) {
	input := &model.Input{
		Zones: []model.ZoneResponse{
			{ZoneId: "z1", Shipments: 3000},
			{ZoneId: "z2", Shipments: 1500},
		},
		Vehicles: []dependency.Vehicle{
			{VehicleID: "v1", MinShipments: 0, MaxShipments: 2300},
			{VehicleID: "v2", MinShipments: 500, MaxShipments: 1500},
		},
		CostPerShipmentMap: map[string]map[string]*dependency.CostPerShipment{
			"v1": {
				"z1": &dependency.CostPerShipment{Cost: 5},
				"z2": &dependency.CostPerShipment{Cost: 7},
			},
			"veh2": {
				"z2": &dependency.CostPerShipment{Cost: 8},
			},
		},
	}

	resp := SolveILP(input)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)
	//assert.Equal(t, uint64(8), resp.Assignments["veh1"]["zone1"]+resp.Assignments["veh2"]["zone1"]+resp.Assignments["veh1"]["zone2"]+resp.Assignments["veh2"]["zone2"])
	assert.Equal(t, uint64(8), resp.TotalCost)
}
