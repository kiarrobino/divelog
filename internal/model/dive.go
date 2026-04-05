package model

import "time"

type Dive struct {
	ID         string
	DiveNumber int
	Date       time.Time
	SiteName   string
	Location   string		// TODO: Coordinates? 
	MaxDepth   float64 		// feet
	AvgDepth   float64		// feet
	Duration   int 			// minutes
	WaterTemp  float64		// fahrenheit
	Visibility int 			// 1-5 
	TankStart  int 			// bar
	TankEnd    int			// bar
	O2Percent  float64 		// 21.0 = air
	WaterType  string  		// "salt" | "fresh" | "brackish"
	DiveType   string  		// "recreational" | "wreck" | "drift" | etc.
	Notes      string
	Rating     int 			// 1–5
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
