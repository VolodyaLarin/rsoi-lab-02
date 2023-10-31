package repo

import (
	"gorm.io/gorm"
	"time"
)

type AirportModel struct {
	gorm.Model

	Name    string
	City    string
	Country string
}

type FlightModel struct {
	gorm.Model

	FlightNumber  string `gorm:"unique"`
	Datetime      time.Time
	FromAirportID int64
	ToAirportID   int64
	Price         int64

	FromAirport *AirportModel `gorm:"foreignKey:CompanyRefer"`
	ToAirport   *AirportModel `gorm:"foreignKey:FromAirportID"`
}
