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
