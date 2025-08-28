package api

import (
	"fmt"
	"shipment_allocation/internal/dependency"

	"gorm.io/gorm"
)

type Allocation struct {
	db *gorm.DB
}

func (a *Allocation) ShipmentAllocation() error {
	vehicle, err := dependency.GetAllVehicles(a.db)
	if err != nil {
		return err
	}

	fmt.Println(vehicle)
	return nil
}

func NewAllocation(db *gorm.DB) *Allocation {
	return &Allocation{db: db}
}
