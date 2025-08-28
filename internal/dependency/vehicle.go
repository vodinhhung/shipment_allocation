package dependency

import "gorm.io/gorm"

type Vehicle struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	MaxShipments uint64 `json:"delivery_latitude"`
	MinShipments uint64 `json:"delivery_longitude"`
}

// GetAllVehicles retrieves all vehicles from the database.
func GetAllVehicles(db *gorm.DB) ([]Vehicle, error) {
	var vehicles []Vehicle
	if err := db.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

// GetVehicleByID retrieves a vehicle by its ID.
func GetVehicleByID(db *gorm.DB, id uint64) (*Vehicle, error) {
	var vehicle Vehicle
	if err := db.First(&vehicle, id).Error; err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (Vehicle) TableName() string {
	return "vehicle_tab"
}
