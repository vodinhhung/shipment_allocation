package api

import (
	"fmt"
	"math"
	"shipment_allocation/internal/dependency"
	"shipment_allocation/internal/model"

	"gorm.io/gorm"
)

type Allocation struct {
	db *gorm.DB
}

func NewAllocation(db *gorm.DB) *Allocation {
	return &Allocation{db: db}
}

func (a *Allocation) ShipmentAllocation(responseData *model.RequestData) error {
	inputVehicleIds := responseData.Vehicles
	vehicleIds := make([]string, 0)
	for _, inputVehicleId := range inputVehicleIds {
		vehicleIds = append(vehicleIds, inputVehicleId)
	}
	fmt.Println(vehicleIds)

	vehicles, err := dependency.GetVehiclesByVehicleIDs(a.db, vehicleIds)
	if err != nil {
		return err
	}

	costPerShipmentMap := map[string]map[string]*dependency.CostPerShipment{}
	for _, zone := range responseData.Zones {
		for _, vehicle := range vehicles {
			costPerShipment, err := dependency.GetCostPerShipmentByZoneAndVehicle(a.db, zone.ZoneId, vehicle.VehicleID)
			if err != nil {
				return err
			}

			if _, ok := costPerShipmentMap[zone.ZoneId]; !ok {
				costPerShipmentMap[vehicle.VehicleID] = map[string]*dependency.CostPerShipment{zone.ZoneId: costPerShipment}
			} else {
				costPerShipmentMap[vehicle.VehicleID][zone.ZoneId] = costPerShipment
			}
		}
	}

	input := &model.Input{
		Zones:              responseData.Zones,
		Vehicles:           vehicles,
		CostPerShipmentMap: costPerShipmentMap,
	}
	response := SolveILP(input)
	fmt.Println(response)

	return nil
}

// ------------------ ILP Solver ------------------

type assignment struct {
	vehicleID string
	zoneID    string
	shipments uint64
}

type solver struct {
	input    *model.Input
	bestCost uint64
	bestPlan []assignment
}

func SolveILP(input *model.Input) *model.Response {
	s := &solver{
		input:    input,
		bestCost: math.MaxUint64,
	}

	assignments := []assignment{}
	used := make(map[string]uint64) // zone -> shipments already allocated
	s.backtrack(0, assignments, 0, used)

	if s.bestCost == math.MaxUint64 {
		// No feasible solution
		return &model.Response{
			Status:     "failed",
			Violations: []string{"No feasible assignment found"},
		}
	}

	// Build assignment map
	result := make(map[string]map[string]uint64)
	for _, a := range s.bestPlan {
		if _, ok := result[a.vehicleID]; !ok {
			result[a.vehicleID] = make(map[string]uint64)
		}
		result[a.vehicleID][a.zoneID] += a.shipments
	}

	return &model.Response{
		Status:      "success",
		Assignments: result,
		TotalCost:   s.bestCost,
	}
}

// Recursive backtracking with pruning
func (s *solver) backtrack(vehicleIdx int, currentPlan []assignment, currentCost uint64, used map[string]uint64) {
	// Base case: all vehicles processed
	if vehicleIdx == len(s.input.Vehicles) {
		// Check if all zones are satisfied
		for _, z := range s.input.Zones {
			if used[z.ZoneId] < z.Shipments {
				return
			}
		}
		// Better solution
		if currentCost < s.bestCost {
			s.bestCost = currentCost
			s.bestPlan = append([]assignment(nil), currentPlan...)
		}
		return
	}

	v := s.input.Vehicles[vehicleIdx]

	// Try all zones this vehicle can serve
	for _, zoneID := range v.Zones {
		zoneDemand := getZoneDemand(s.input.Zones, zoneID)
		if zoneDemand == 0 {
			continue
		}

		maxCanTake := minU64(v.MaxShipments, zoneDemand-used[zoneID])

		// Try quantities from MinShipments → maxCanTake
		for qty := v.MinShipments; qty <= maxCanTake; qty++ {
			cost := uint64(0)
			if cps, ok := s.input.CostPerShipmentMap[v.VehicleID][zoneID]; ok {
				cost = cps.Cost * qty
			}

			used[zoneID] += qty
			newPlan := append(currentPlan, assignment{v.VehicleID, zoneID, qty})
			newCost := currentCost + cost

			// Prune if already worse than best
			if newCost < s.bestCost {
				s.backtrack(vehicleIdx+1, newPlan, newCost, used)
			}

			// Backtrack
			used[zoneID] -= qty
		}
	}

	// Option: vehicle unused if MinShipments == 0
	if v.MinShipments == 0 {
		s.backtrack(vehicleIdx+1, currentPlan, currentCost, used)
	}
}

// ------------------ Helpers ------------------
func getZoneDemand(zones []model.ZoneResponse, id string) uint64 {
	for _, z := range zones {
		if z.ZoneId == id {
			return z.Shipments
		}
	}
	return 0
}

func minU64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
