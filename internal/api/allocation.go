package api

import (
	"fmt"
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

func (a *Allocation) ShipmentAllocation(responseData *model.RequestData) (*model.Response, error) {
	inputVehicleIds := responseData.Vehicles
	vehicleIds := make([]string, 0)
	for _, inputVehicleId := range inputVehicleIds {
		vehicleIds = append(vehicleIds, inputVehicleId)
	}
	fmt.Println(vehicleIds)

	vehicles, err := dependency.GetVehiclesByVehicleIDs(a.db, vehicleIds)
	if err != nil {
		return nil, err
	}

	costPerShipmentMap := map[string]map[string]*model.CostPerShipment{}
	for _, zone := range responseData.Zones {
		for _, vehicle := range vehicles {
			costPerShipment, err := dependency.GetCostPerShipmentByZoneAndVehicle(a.db, zone.ZoneId, vehicle.VehicleID)
			if err != nil {
				continue
			}

			if _, ok := costPerShipmentMap[vehicle.VehicleID]; !ok {
				costPerShipmentMap[vehicle.VehicleID] = map[string]*model.CostPerShipment{zone.ZoneId: costPerShipment}
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

	return response, nil
}
