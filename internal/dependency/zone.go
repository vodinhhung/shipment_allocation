package dependency

import "gorm.io/gorm"

type Zone struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	NoShipments uint64 `json:"no_shipments"`
}

// GetAllZones retrieves all zones from the database.
func GetAllZones(db *gorm.DB) ([]Zone, error) {
	var zones []Zone
	if err := db.Find(&zones).Error; err != nil {
		return nil, err
	}
	return zones, nil
}

// GetZoneByID retrieves a zone by its ID.
func GetZoneByID(db *gorm.DB, id uint64) (*Zone, error) {
	var zone Zone
	if err := db.First(&zone, id).Error; err != nil {
		return nil, err
	}
	return &zone, nil
}
