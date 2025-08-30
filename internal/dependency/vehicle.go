package dependency

import (
	"shipment_allocation/internal/model"

	"gorm.io/gorm"
)

// GetAllVehicles retrieves all vehicles from the database.
func GetAllVehicles(db *gorm.DB) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	if err := db.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

// GetVehicleByID retrieves a vehicle by its ID.
func GetVehicleByID(db *gorm.DB, id uint64) (*model.Vehicle, error) {
	var vehicle model.Vehicle
	if err := db.First(&vehicle, id).Error; err != nil {
		return nil, err
	}
	return &vehicle, nil
}

// GetVehiclesByIDs retrieves vehicles by a list of IDs.
func GetVehiclesByIDs(db *gorm.DB, ids []uint64) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	if err := db.Where("id IN ?", ids).Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

// GetVehiclesByVehicleIDs retrieves vehicles by a list of vehicle_id strings.
func GetVehiclesByVehicleIDs(db *gorm.DB, vehicleIDs []string) ([]model.Vehicle, error) {
	var vehicles []model.Vehicle
	if err := db.Where("vehicle_id IN ?", vehicleIDs).Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}
