package repo

import (
	"context"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/ticket"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormTicketRepo struct {
	db *gorm.DB
}

func NewTicketRepo(db *gorm.DB) *GormTicketRepo {
	return &GormTicketRepo{
		db: db,
	}
}

func modelToDto(model TicketModel) bonus.TicketDto {
	return bonus.TicketDto{
		ID:           int64(model.ID),
		Username:     model.Username,
		TicketUid:    model.TicketUid,
		FlightNumber: model.FlightNumber,
		Price:        int(model.Price),
		Status:       string(model.Status),
	}
}

func (g GormTicketRepo) List(ctx context.Context, filter bonus.TicketFilter) (error, []bonus.TicketDto) {
	var tickets []TicketModel

	query := g.db.Model(&TicketModel{
		Username: filter.Username,
	})
	if len(filter.Uids) != 0 {
		query = query.Where("ticket_uuid in ?", filter.Uids)
	}

	err := query.Find(&tickets).Error
	if err != nil {
		return err, nil
	}

	return nil, utils.Map(tickets, modelToDto)
}

func (g GormTicketRepo) FindTicketByUid(ctx context.Context, uuid uuid.UUID) (error, *bonus.TicketDto) {
	ticket := TicketModel{}
	err := g.db.Model(TicketModel{TicketUid: uuid}).First(&ticket).Error
	if err != nil {
		return err, nil
	}
	dto := modelToDto(ticket)
	return nil, &dto

}

func (g GormTicketRepo) SaveTicket(ctx context.Context, dto bonus.TicketDto) (error, *bonus.TicketDto) {
	ticket := TicketModel{
		Model: gorm.Model{
			ID: uint(dto.ID),
		},
		TicketUid:    dto.TicketUid,
		Username:     dto.Username,
		FlightNumber: dto.FlightNumber,
		Price:        int64(dto.Price),
		Status:       TicketStatus(dto.Status),
	}
	err := g.db.Save(&ticket).Error
	if err != nil {
		return err, nil
	}
	resDto := modelToDto(ticket)
	return nil, &resDto

}
