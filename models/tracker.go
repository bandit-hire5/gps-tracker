package models

import "time"

type Tracker struct {
	IMEI    string
	Command string
	Date    time.Time
	Lat     float64
	Lon     float64
	Speed   float64
}
