package usecase

import (
	"context"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus"
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus/mock_bonus"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestBonusUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_bonus.NewMockIBonusRepo(ctrl)

	uc := NewBonusUsecase(repo)

	dto := bonus.BonusDto{
		Balance: 500,
		Status:  bonus.PrivelegeSilver,
		History: nil,
	}

	repo.EXPECT().GetBonusDetails(gomock.Any(), gomock.Eq("testusername")).Return(nil, &dto).Times(1)

	err, res := uc.List(context.Background(), "testusername")

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.Balance, int64(500))
	assert.Equal(t, res.Status, bonus.PrivelegeSilver)
	assert.Nil(t, res.History)

}

func TestBonusUsecase_MinusFromPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_bonus.NewMockIBonusRepo(ctrl)

	uc := NewBonusUsecase(repo)

	dto := bonus.BonusDto{
		Balance: 300,
		Status:  bonus.PrivelegeSilver,
		History: nil,
	}
	repo.EXPECT().GetBonusDetails(gomock.Any(), gomock.Eq("testusername")).Return(nil, &dto).Times(1)

	repo.EXPECT().CreateBonusOperation(gomock.Any(), gomock.Eq("testusername"), gomock.Cond(func(x any) bool {
		dto, ok := x.(bonus.BonusHistoryDto)
		if !ok {
			return false
		}

		return dto.OperationType == bonus.PrivilegeHistoryDebit && dto.BalanceDiff == -300
	})).Return(nil, nil, nil).Times(1)

	err, item, b := uc.MinusFromPayment(context.Background(), "testusername", uuid.UUID{}, 400, false)

	assert.Nil(t, err)
	assert.Nil(t, item)
	assert.Nil(t, b)
}

func TestBonusUsecase_PlusFromPayment(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock_bonus.NewMockIBonusRepo(ctrl)

	uc := NewBonusUsecase(repo)

	repo.EXPECT().CreateBonusOperation(gomock.Any(), gomock.Eq("testusername"), gomock.Cond(func(x any) bool {
		dto, ok := x.(bonus.BonusHistoryDto)
		if !ok {
			return false
		}

		return dto.OperationType == bonus.PrivilegeHistoryFill && dto.BalanceDiff == 40
	})).Return(nil, nil, nil).Times(1)

	err, item, b := uc.PlusFromPayment(context.Background(), "testusername", uuid.UUID{}, 400, false)

	assert.Nil(t, err)
	assert.Nil(t, item)
	assert.Nil(t, b)
}
