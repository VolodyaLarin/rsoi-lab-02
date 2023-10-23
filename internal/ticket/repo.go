package bonus

import (
	"context"
	"github.com/google/uuid"
)

type ITicketRepo interface {
	List(ctx context.Context, filter TicketFilter) (error, []TicketDto)
	FindTicketByUid(ctx context.Context, uuid uuid.UUID) (error, *TicketDto)
	SaveTicket(ctx context.Context, dto TicketDto) (error, *TicketDto)
}
