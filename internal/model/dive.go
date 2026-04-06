package model

import "time"

type Dive struct {
	ID         string
	DiveNumber int
	Date       time.Time `json:"date"`
	SiteName   string    `json:"site_name"`
	Location   string    `json:"location"`   // TODO: Coordinates?
	MaxDepth   float64   `json:"max_depth"`  // feet
	AvgDepth   float64   `json:"avg_depth"`  // feet
	Duration   int       `json:"duration"`   // minutes
	WaterTemp  float64   `json:"water_temp"` // fahrenheit
	Visibility float64   `json:"visibility"` // feet
	TankStart  int       `json:"tank_start"` // bar
	TankEnd    int       `json:"tank_end"`   // bar
	O2Percent  float64   `json:"o2_percent"` // 21.0 = air
	WaterType  string    `json:"water_type"` // "salt" | "fresh" | "brackish"
	DiveType   string    `json:"dive_type"`  // "recreational" | "wreck" | "drift" | etc.
	Notes      string    `json:"notes"`
	Rating     int       `json:"rating"` // 1–5
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateDiveInput struct {
	Date       string  `json:"date"`
	SiteName   string  `json:"site_name"`
	Location   string  `json:"location"`
	MaxDepth   float64 `json:"max_depth"`
	AvgDepth   float64 `json:"avg_depth"`
	Duration   int     `json:"duration"`
	WaterTemp  float64 `json:"water_temp"`
	Visibility float64 `json:"visibility"`
	TankStart  int     `json:"tank_start"`
	TankEnd    int     `json:"tank_end"`
	O2Percent  float64 `json:"o2_percent"`
	WaterType  string  `json:"water_type"`
	DiveType   string  `json:"dive_type"`
	Notes      string  `json:"notes"`
	Rating     int     `json:"rating"`
}
