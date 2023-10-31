package repo

import (
	"github.com/VolodyaLarin/rsoi-lab-02/internal/bonus"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type PrivilegeHistoryModel struct {
	gorm.Model

	PrivilegeID   int64
	TicketUid     uuid.UUID
	Datetime      time.Time
	BalanceDiff   int64
	OperationType bonus.PrivilegeHistoryOperationType
}

type PrivilegeModel struct {
	gorm.Model

	ID       int64
	Username string
	Status   bonus.PrivilegeStatus
	Balance  int64

	History []PrivilegeHistoryModel `gorm:"foreignKey:PrivilegeID"`
}
