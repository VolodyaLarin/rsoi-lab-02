package bonus

import "github.com/google/uuid"

type TicketDto struct {
	ID           int64     `json:"-"`
	Username     string    `json:"-"`
	TicketUid    uuid.UUID `json:"ticketUid"`
	FlightNumber string    `json:"flightNumber"`
	Price        int       `json:"price"`
	Status       string    `json:"status"`
}
type TicketFilter struct {
	Username string
	Uids     []uuid.UUID
}
