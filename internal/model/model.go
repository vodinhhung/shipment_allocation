package model

type RequestData struct {
	Zones    []ZoneResponse `json:"zones"`
	Vehicles []string       `json:"vehicles"`
}

type ZoneResponse struct {
	ZoneId    string `json:"zone_id"`
	Shipments uint64 `json:"shipments"`
}

type Input struct {
	Zones              []ZoneResponse `json:"zones"`
	Vehicles           []Vehicle
	CostPerShipmentMap map[string]map[string]*CostPerShipment
}

type Response struct {
	Status      string                       `json:"status"` // "success" | "failed"
	Assignments map[string]map[string]uint64 `json:"assignments,omitempty"`
	TotalCost   uint64                       `json:"total_cost,omitempty"`
}

type CostPerShipment struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	VehicleID string `json:"vehicle_id"`
	ZoneID    string `json:"zone_id"`
	Cost      uint64 `json:"cost"`
}

func (CostPerShipment) TableName() string {
	return "cost_per_shipment_tab"
}

type Vehicle struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	VehicleID    string `json:"vehicle_id"`
	MaxShipments uint64 `json:"delivery_latitude"`
	MinShipments uint64 `json:"delivery_longitude"`
}

func (Vehicle) TableName() string {
	return "vehicle_tab"
}
