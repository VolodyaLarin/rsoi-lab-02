package bonus

import (
	"context"
	"github.com/google/uuid"
)

type IBonusRepo interface {
	GetBonusDetails(ctx context.Context, username string) (error, *BonusDto)

	CreateBonusOperation(ctx context.Context, username string, dto BonusHistoryDto) (error, *BonusHistoryDto, *BonusDto)
	DeleteBonusOperationByFlightUid(ctx context.Context, username string, ticketUuid uuid.UUID) (error, *BonusHistoryDto)
}
