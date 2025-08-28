package dependency

import "gorm.io/gorm"

type CostPerShipment struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	VehicleID uint64 `json:"vehicle_id"`
	ZoneID    uint64 `json:"zone_id"`
	Cost      uint64 `json:"cost"`
}

func GetAllCostPerShipments(db *gorm.DB) ([]CostPerShipment, error) {
	var costs []CostPerShipment
	err := db.Find(&costs).Error
	return costs, err
}

func GetCostPerShipmentsByZoneID(db *gorm.DB, zoneID uint64) ([]CostPerShipment, error) {
	var costs []CostPerShipment
	err := db.Where("zone_id = ?", zoneID).Find(&costs).Error
	return costs, err
}

func GetCostPerShipmentsByVehicleID(db *gorm.DB, vehicleID uint64) ([]CostPerShipment, error) {
	var costs []CostPerShipment
	err := db.Where("vehicle_id = ?", vehicleID).Find(&costs).Error
	return costs, err
}

func GetCostPerShipmentByZoneAndVehicle(db *gorm.DB, zoneID, vehicleID uint64) (*CostPerShipment, error) {
	var cost CostPerShipment
	err := db.Where("zone_id = ? AND vehicle_id = ?", zoneID, vehicleID).First(&cost).Error
	return &cost, err
}

func (CostPerShipment) TableName() string {
	return "cost_per_shipment_tab"
}
