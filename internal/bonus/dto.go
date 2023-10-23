package bonus

import (
	"github.com/google/uuid"
	"time"
)

type PrivilegeHistoryOperationType string

const (
	PrivilegeHistoryFill  PrivilegeHistoryOperationType = "FILL_IN_BALANCE"
	PrivilegeHistoryDebit PrivilegeHistoryOperationType = "DEBIT_THE_ACCOUNT"
)

type PrivilegeStatus string

const (
	PrivelegeBronze PrivilegeStatus = "BRONZE"
	PrivelegeSilver PrivilegeStatus = "SILVER"
	PrivelegeGold   PrivilegeStatus = "GOLD"
)

type BonusHistoryDto struct {
	Date          time.Time                     `json:"date"`
	TicketUid     uuid.UUID                     `json:"ticketUid"`
	BalanceDiff   int64                         `json:"balanceDiff"`
	OperationType PrivilegeHistoryOperationType `json:"operationType"`
}
type BonusDto struct {
	Balance int64             `json:"balance"`
	Status  PrivilegeStatus   `json:"status"`
	History []BonusHistoryDto `json:"history"`
}
