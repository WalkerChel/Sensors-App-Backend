package apiResponses

type SensorWithoutRegionResponse struct {
	ID        int64   `json:"sensor_id"`
	Name      string  `json:"sensor_name"`
	Longitude float64 `json:"sensor_longitude"`
	Latitude  float64 `json:"sensor_latitude"`
}

type SensorsInRegionResponse struct {
	Sensors  []SensorWithoutRegionResponse `json:"sensors"`
	RegionID int64                         `json:"region_id"`
}

type SensorsPaginatedResponse struct {
	Sensors         []SensorWithoutRegionResponse `json:"sensors"`
	AllSensorsCount int64                         `json:"all_sensors_amount,omitempty"`
}
