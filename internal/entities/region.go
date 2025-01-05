package entities

type Region struct {
	Id   int64  `json:"region_id" db:"id"`
	Name string `json:"region_name" db:"name"`
}
