package repo

import (
	"context"
	"errors"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormBonusRepo struct {
	db *gorm.DB
}

func NewBonusRepo(db *gorm.DB) *GormBonusRepo {
	return &GormBonusRepo{db: db}
}
func (g GormBonusRepo) getBonuses(ctx context.Context, username string) (error, PrivilegeModel) {
	model := PrivilegeModel{}
	err := g.db.Preload("History").Model(PrivilegeModel{Username: username}).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		model = PrivilegeModel{
			Username: username,
			Status:   bonus.PrivelegeBronze,
			Balance:  0,
		}
		err = nil
	}

	if err != nil {
		return err, PrivilegeModel{}
	}

	return nil, model
}

func (g GormBonusRepo) GetBonusDetails(ctx context.Context, username string) (error, *bonus.BonusDto) {
	err, model := g.getBonuses(ctx, username)
	if err != nil {
		return err, nil
	}

	return nil, &bonus.BonusDto{
		Balance: model.Balance,
		Status:  model.Status,
		History: utils.Map(model.History, func(item PrivilegeHistoryModel) bonus.BonusHistoryDto {
			return bonus.BonusHistoryDto{
				Date:          item.Datetime,
				TicketUid:     item.TicketUid,
				BalanceDiff:   item.BalanceDiff,
				OperationType: item.OperationType,
			}
		}),
	}
}

func (g GormBonusRepo) CreateBonusOperation(ctx context.Context, username string, dto bonus.BonusHistoryDto) (error, *bonus.BonusHistoryDto, *bonus.BonusDto) {

	err, model := g.getBonuses(ctx, username)
	if err != nil {
		return err, nil, nil
	}

	history := PrivilegeHistoryModel{
		PrivilegeID:   model.ID,
		TicketUid:     dto.TicketUid,
		Datetime:      dto.Date,
		BalanceDiff:   dto.BalanceDiff,
		OperationType: dto.OperationType,
	}

	model.Balance += dto.BalanceDiff

	err = g.db.Transaction(func(tx *gorm.DB) error {
		err = g.db.Save(&model).Error
		if err != nil {
			return err
		}

		err = g.db.Save(&history).Error
		if err != nil {
			return err
		}

		return nil
	})

	return nil, &bonus.BonusHistoryDto{
			Date:          history.Datetime,
			TicketUid:     history.TicketUid,
			BalanceDiff:   history.BalanceDiff,
			OperationType: history.OperationType,
		}, &bonus.BonusDto{
			Balance: model.Balance,
			Status:  model.Status,
		}
}

func (g GormBonusRepo) DeleteBonusOperationByFlightUid(ctx context.Context, username string, ticketUuid uuid.UUID) (error, *bonus.BonusHistoryDto) {
	err, model := g.getBonuses(ctx, username)
	if err != nil {
		return err, nil
	}

	item := PrivilegeHistoryModel{}
	err = g.db.Model(PrivilegeHistoryModel{TicketUid: ticketUuid}).First(&item).Error
	if err != nil {
		return err, nil
	}

	model.Balance -= item.BalanceDiff
	if model.Balance < 0 {
		model.Balance = 0
	}

	err = g.db.Transaction(func(tx *gorm.DB) error {
		err = g.db.Save(&model).Error
		if err != nil {
			return err
		}

		err = g.db.Delete(&PrivilegeHistoryModel{}, item.ID).Error
		if err != nil {
			return err
		}

		return nil
	})

	return nil, &bonus.BonusHistoryDto{
		Date:          item.Datetime,
		TicketUid:     item.TicketUid,
		BalanceDiff:   item.BalanceDiff,
		OperationType: item.OperationType,
	}

}
