package usecase

import (
	"context"
	"errors"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket/repo"
	"github.com/google/uuid"
)

type TicketUsecase struct {
	repo bonus.ITicketRepo
}

func NewTicketUsecase(repo bonus.ITicketRepo) *TicketUsecase {
	return &TicketUsecase{repo: repo}
}

func (uc TicketUsecase) List(ctx context.Context, filter *bonus.TicketFilter) (error, []bonus.TicketDto) {
	return uc.repo.List(ctx, *filter)
}

func (uc TicketUsecase) Buy(ctx context.Context, username string, flightNumber string, price int64) (error, *bonus.TicketDto) {
	return uc.repo.SaveTicket(ctx, bonus.TicketDto{
		TicketUid:    uuid.New(),
		Username:     username,
		FlightNumber: flightNumber,
		Price:        int(price),
		Status:       string(repo.TicketPaid),
	})
}

func (uc TicketUsecase) Cancel(ctx context.Context, username string, uid uuid.UUID) (error, *bonus.TicketDto) {
	err, ticket := uc.repo.FindTicketByUid(ctx, uid)
	if err != nil {
		return err, nil
	}
	if ticket.Username != username {
		return errors.New("forbidden"), nil
	}
	if ticket.Status != string(repo.TicketPaid) {
		return errors.New("operation not supported"), nil
	}

	ticket.Status = string(repo.TicketCanceled)
	return uc.repo.SaveTicket(ctx, *ticket)
}
