package dependency

import (
	"shipment_allocation/internal/model"

	"gorm.io/gorm"
)

func GetAllCostPerShipments(db *gorm.DB) ([]model.CostPerShipment, error) {
	var costs []model.CostPerShipment
	err := db.Find(&costs).Error
	return costs, err
}

func GetCostPerShipmentsByZoneID(db *gorm.DB, zoneID string) ([]model.CostPerShipment, error) {
	var costs []model.CostPerShipment
	err := db.Where("zone_id = ?", zoneID).Find(&costs).Error
	return costs, err
}

func GetCostPerShipmentsByVehicleID(db *gorm.DB, vehicleID string) ([]model.CostPerShipment, error) {
	var costs []model.CostPerShipment
	err := db.Where("vehicle_id = ?", vehicleID).Find(&costs).Error
	return costs, err
}

func GetCostPerShipmentByZoneAndVehicle(db *gorm.DB, zoneID, vehicleID string) (*model.CostPerShipment, error) {
	var cost model.CostPerShipment
	err := db.Where("zone_id = ? AND vehicle_id = ?", zoneID, vehicleID).First(&cost).Error
	return &cost, err
}
