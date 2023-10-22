package repo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketStatus string

const (
	TicketPaid     TicketStatus = "PAID"
	TicketCanceled TicketStatus = "CANCELED"
)

type TicketModel struct {
	gorm.Model

	TicketUid    uuid.UUID `gorm:"unique"`
	Username     string
	FlightNumber string
	Price        int64
	Status       TicketStatus
}
