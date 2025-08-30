package api

import (
	"fmt"
	"math"
	"shipment_allocation/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSolveILP_BasicCase(t *testing.T) {
	input := &model.Input{
		Zones: []model.ZoneResponse{
			{ZoneId: "z1", Shipments: 2000},
			{ZoneId: "z2", Shipments: 1500},
		},
		Vehicles: []model.Vehicle{
			{VehicleID: "v1", MinShipments: 0, MaxShipments: 2300},
			{VehicleID: "v2", MinShipments: 500, MaxShipments: 1500},
		},
		CostPerShipmentMap: map[string]map[string]*model.CostPerShipment{
			"v1": {
				"z1": &model.CostPerShipment{Cost: 5},
				"z2": &model.CostPerShipment{Cost: 7},
			},
			"v2": {
				"z2": &model.CostPerShipment{Cost: 8},
			},
		},
	}
	/*
		v1: z1: 2300, z2: 7
		v2: z2: 8
		total cost: 18700
	*/

	resp := SolveILP(input)
	assert.NotNil(t, resp)
	assert.Equal(t, "success", resp.Status)
	assert.Equal(t, uint64(18700), resp.TotalCost)
}

func TestBacktrackAndAssignZone(t *testing.T) {
	// Create test input
	input := &model.Input{
		Zones: []model.ZoneResponse{
			{ZoneId: "z1", Shipments: 2000},
			{ZoneId: "z2", Shipments: 1500},
		},
		Vehicles: []model.Vehicle{
			{VehicleID: "v1", MinShipments: 0, MaxShipments: 2300},
			{VehicleID: "v2", MinShipments: 500, MaxShipments: 1500},
		},
		CostPerShipmentMap: map[string]map[string]*model.CostPerShipment{
			"v1": {
				"z1": &model.CostPerShipment{Cost: 5},
				"z2": &model.CostPerShipment{Cost: 7},
			},
			"v2": {
				"z2": &model.CostPerShipment{Cost: 8},
			},
		},
	}

	// Create solver instance
	s := &solver{
		input:    input,
		bestCost: math.MaxUint64,
		bestPlan: []assignment{},
	}

	// Setup for backtrackZones test
	vehicleUsed := make(map[string]uint64)
	current := make(map[string]map[string]uint64)
	for _, v := range input.Vehicles {
		vehicleUsed[v.VehicleID] = 0
		current[v.VehicleID] = make(map[string]uint64)
	}

	bestCost := uint64(math.MaxUint64)
	bestAssign := make(map[string]map[string]uint64)

	// Test backtrackZones starting from zone 0
	s.backtrackZones(
		input,
		0, // start with zone index 0
		vehicleUsed,
		current,
		0, // current cost = 0
		&bestCost,
		&bestAssign,
	)

	// Expected optimal assignments:
	// v1: 200 from z1 (cost 5*200=1000), 0 from z2
	// v2: 0 from z1, 300 from z2 (cost 7*300=2100)
	// Total cost: 3100

	assert.Equal(t, uint64(21700), bestCost, "Expected optimal cost of 21700")
	fmt.Println(fmt.Sprintf("Best cost: %d", bestCost))

	expectedAssign := map[string]map[string]uint64{
		"v1": {
			"z1": 2000,
			"z2": 300,
		},
		"v2": {
			"z2": 1200,
		},
	}
	assert.Equal(t, expectedAssign, bestAssign, "Expected optimal cost of 21700")
	fmt.Println(fmt.Sprintf("Best assign: %v", bestAssign))
}
