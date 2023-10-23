package usecase

import (
	"context"
	"time"
)

type FlightDto struct {
	ID           int64     `json:"-"`
	FlightNumber string    `json:"flightNumber"`
	FromAirport  string    `json:"fromAirport"`
	ToAirport    string    `json:"toAirport"`
	Date         time.Time `json:"date"`
	Price        int       `json:"price"`
}

type FlightUsecase struct {
}

type FlightFilter struct {
	Offset   int64
	Limit    int64
	Username string
	Numbers  []string
}

func (uc FlightUsecase) List(ctx context.Context, filter *FlightFilter) (error, []FlightDto, int64) {
	return nil, []FlightDto{
		{
			FlightNumber: "AFL031",
			FromAirport:  "Санкт-Петербург Пулково",
			ToAirport:    "Москва Шереметьево",
			Date:         time.Date(2021, 10, 8, 20, 0, 0, 0, time.Local),
			Price:        1500,
		},
		{
			FlightNumber: "aaasd-dsd",
			FromAirport:  "XX",
			ToAirport:    "zz",
			Date:         time.Now(),
			Price:        30000,
		},
	}, 2
}
