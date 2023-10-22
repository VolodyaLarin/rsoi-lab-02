package usecase

import (
	"context"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus"
	"github.com/google/uuid"
	"math"
	"time"
)

type BonusUsecase struct {
	repo bonus.IBonusRepo
}

func NewBonusUsecase(repo bonus.IBonusRepo) *BonusUsecase {
	return &BonusUsecase{repo: repo}
}

func (uc BonusUsecase) List(ctx context.Context, username string) (error, *bonus.BonusDto) {
	return uc.repo.GetBonusDetails(ctx, username)
}

func (uc BonusUsecase) PlusFromPayment(ctx context.Context, username string, ticketUid uuid.UUID, payment int64, draft bool) (error, *bonus.BonusHistoryDto, *bonus.BonusDto) {
	balanceDiff := payment / 10

	err, history, dto := uc.repo.CreateBonusOperation(ctx, username, bonus.BonusHistoryDto{
		Date:          time.Now(),
		TicketUid:     ticketUid,
		BalanceDiff:   balanceDiff,
		OperationType: bonus.PrivilegeHistoryFill,
	})
	if err != nil {
		return err, nil, nil
	}

	return nil, history, dto

}
func (uc BonusUsecase) MinusFromPayment(ctx context.Context, username string, ticketUid uuid.UUID, payment int64, draft bool) (error, *bonus.BonusHistoryDto, *bonus.BonusDto) {
	balanceDiff := -int64(math.Min(200, float64(payment)))

	err, history, dto := uc.repo.CreateBonusOperation(ctx, username, bonus.BonusHistoryDto{
		Date:          time.Now(),
		TicketUid:     ticketUid,
		BalanceDiff:   balanceDiff,
		OperationType: bonus.PrivilegeHistoryDebit,
	})

	if err != nil {
		return err, nil, nil
	}

	return nil, history, dto
}

func (uc BonusUsecase) RevertFromPayment(ctx context.Context, username string, ticketUid uuid.UUID) error {
	err, _ := uc.repo.DeleteBonusOperationByFlightUid(ctx, username, ticketUid)
	if err != nil {
		return err
	}
	return nil
}
