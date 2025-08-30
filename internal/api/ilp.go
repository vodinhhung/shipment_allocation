package api

import (
	"math"
	"shipment_allocation/internal/model"
)

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
		input: input,
	}

	bestCost := uint64(math.MaxUint64)
	bestAssign := make(map[string]map[string]uint64)

	// remaining demand for each zone
	remaining := make(map[string]uint64)
	for _, z := range input.Zones {
		remaining[z.ZoneId] = z.Shipments
	}

	currentAssign := make(map[string]map[string]uint64)

	s.backtrackZones(input, 0, remaining, currentAssign, 0, &bestCost, &bestAssign)

	if bestCost == math.MaxUint64 {
		return &model.Response{
			Status: "failed",
		}
	}

	return &model.Response{
		Status:      "success",
		Assignments: bestAssign,
		TotalCost:   bestCost,
	}
}

func (s *solver) assignZone(
	input *model.Input,
	zoneIdx int,
	vehIdx int,
	need uint64,
	vehicleUsed map[string]uint64,
	current map[string]map[string]uint64,
	cost uint64,
	bestCost *uint64,
	bestAssign *map[string]map[string]uint64,
) {
	// prune by cost
	if cost >= *bestCost {
		return
	}

	// all vehicles considered for this zone
	if vehIdx >= len(input.Vehicles) {
		if need == 0 {
			// zone fully assigned -> move to next zone
			s.backtrackZones(input, zoneIdx+1, vehicleUsed, current, cost, bestCost, bestAssign)
		}
		// else fail (not all demand served by available vehicles)
		return
	}

	v := input.Vehicles[vehIdx]
	vid := v.VehicleID

	// compute how much this vehicle can still take (global remaining capacity)
	remainCap := uint64(0)
	if v.MaxShipments > vehicleUsed[vid] {
		remainCap = v.MaxShipments - vehicleUsed[vid]
	} else {
		remainCap = 0
	}

	// If vehicle does not serve this zone => skip it
	cpsRow, ok := input.CostPerShipmentMap[vid]
	if !ok {
		// skip this vehicle
		s.assignZone(input, zoneIdx, vehIdx+1, need, vehicleUsed, current, cost, bestCost, bestAssign)
		return
	}
	cps, serves := cpsRow[input.Zones[zoneIdx].ZoneId]
	if !serves || remainCap == 0 {
		// either does not serve this zone, or no capacity left
		s.assignZone(input, zoneIdx, vehIdx+1, need, vehicleUsed, current, cost, bestCost, bestAssign)
		return
	}

	// upper bound this vehicle can take for this zone
	maxTake := need
	if remainCap < maxTake {
		maxTake = remainCap
	}

	// Additional pruning: check sum of remaining capacities of vehicles[vehIdx..] >= need
	var totalAvail uint64
	for j := vehIdx; j < len(input.Vehicles); j++ {
		vv := input.Vehicles[j]
		rem := uint64(0)
		if vv.MaxShipments > vehicleUsed[vv.VehicleID] {
			rem = vv.MaxShipments - vehicleUsed[vv.VehicleID]
		}
		// but if vv doesn't serve zone, rem contributes 0
		if row, ok := input.CostPerShipmentMap[vv.VehicleID]; ok {
			if _, ok2 := row[input.Zones[zoneIdx].ZoneId]; ok2 {
				totalAvail += rem
			}
		}
	}
	if totalAvail < need {
		// cannot satisfy this zone from remaining vehicles => prune
		return
	}

	// Try options: how much this vehicle will take from this zone (0..maxTake)
	// (We include 0 to allow skipping this vehicle for this zone.)
	for take := uint64(0); take <= maxTake; take++ {
		// apply
		if take > 0 {
			vehicleUsed[vid] += take
			if current[vid] == nil {
				current[vid] = make(map[string]uint64)
			}
			current[vid][input.Zones[zoneIdx].ZoneId] += take
		}

		addedCost := take * cps.Cost
		nextNeed := need - take

		// recurse
		s.assignZone(input, zoneIdx, vehIdx+1, nextNeed, vehicleUsed, current, cost+addedCost, bestCost, bestAssign)

		// revert
		if take > 0 {
			current[vid][input.Zones[zoneIdx].ZoneId] -= take
			if current[vid][input.Zones[zoneIdx].ZoneId] == 0 {
				delete(current[vid], input.Zones[zoneIdx].ZoneId)
			}
			vehicleUsed[vid] -= take
		}
	}
}

// backtrackZones: assign zones sequentially
func (s *solver) backtrackZones(
	input *model.Input,
	zoneIdx int,
	vehicleUsed map[string]uint64,
	current map[string]map[string]uint64,
	cost uint64,
	bestCost *uint64,
	bestAssign *map[string]map[string]uint64,
) {
	// prune by cost
	if cost >= *bestCost {
		return
	}

	// if all zones done -> final check: vehicle mins
	if zoneIdx >= len(input.Zones) {
		// ensure each vehicle meets MinShipments
		for _, v := range input.Vehicles {
			if vehicleUsed[v.VehicleID] < v.MinShipments {
				return
			}
		}
		// feasible solution
		*bestCost = cost
		*bestAssign = deepCopy(current)
		return
	}

	// assign this zone demand using vehicles starting from vehicle 0
	need := input.Zones[zoneIdx].Shipments
	s.assignZone(input, zoneIdx, 0, need, vehicleUsed, current, cost, bestCost, bestAssign)
}

// deepCopy helper for map[string]map[string]uint64
func deepCopy(src map[string]map[string]uint64) map[string]map[string]uint64 {
	cop := make(map[string]map[string]uint64)
	for v, zones := range src {
		cop[v] = make(map[string]uint64)
		for z, val := range zones {
			cop[v][z] = val
		}
	}
	return cop
}
