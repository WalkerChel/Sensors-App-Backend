package entities

type DateTime struct {
	Year   int `json:"year" binding:"required"`
	Month  int `json:"month" binding:"required"`
	Day    int `json:"day" binding:"required"`
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
