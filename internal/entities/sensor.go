package entities

type Sensor struct {
	Id        int64   `json:"sensor_id" db:"id"`
	RegionId  int64   `json:"region_id" db:"region_id"`
	Name      string  `json:"sensor_name" db:"name"`
	Longitude float64 `json:"sensor_longitude" db:"longitude"`
	Latitude  float64 `json:"sensor_latitude" db:"latitude"`
}
