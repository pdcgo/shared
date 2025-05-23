package db_models

import (
	"time"

	"github.com/pdcgo/shared/interfaces/identity_iface"
)

type Withdrawal struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	CreatedByID   uint      `json:"created_by_id"`
	HistID        uint      `json:"hist_id"`
	TeamID        uint      `json:"team_id"`
	MpID          uint      `json:"mp_id"`
	OrderNotFound int       `json:"order_not_found"`
	OrderValid    int       `json:"order_valid"`
	DiffAmount    float64   `json:"diff_amount"`
	AfterAmount   float64   `json:"after_amount"`
	At            time.Time `json:"at" gorm:"index"`
	IsNew         bool      `json:"-" gorm:"-"`

	Hist      *AssetHistory `json:"hist"`
	CreatedBy *User         `json:"created_by"`
}

type OrderAdjustment struct {
	ID      uint `json:"id" gorm:"primarykey"`
	OrderID uint `json:"order_id"`
	MpID    uint `json:"mp_id"`

	At      time.Time      `json:"at" gorm:"index"`
	FundAt  time.Time      `json:"fund_at" gorm:"index"`
	Type    AdjustmentType `json:"type"`
	Amount  float64        `json:"amount"`
	Desc    string         `json:"desc"`
	Deleted bool           `json:"deleted" gorm:"index"`

	Order *Order       `json:"order"`
	Mp    *Marketplace `json:"mp"`
}

func (OrderAdjustment) TableName() string {
	return "order_adjustments"
}

type OrderAdjLogType string

const (
	AdjLogCreated OrderAdjLogType = "created"
	AdjLogUpdated OrderAdjLogType = "updated"
	AdjLogDeleted OrderAdjLogType = "deleted"
)

type OrderAdjustmentLog struct {
	ID        uint                     `json:"id" gorm:"primarykey"`
	AdjID     uint                     `json:"adj_id"`
	UserID    uint                     `json:"user_id"`
	From      identity_iface.AgentType `json:"from"`
	LogType   OrderAdjLogType          `json:"log_type"`
	Timestamp time.Time                `json:"timestamp" gorm:"index"`

	Adj  *OrderAdjustment `json:"-"`
	User *User            `json:"-"`
}

type WdValid struct {
	ID                uint `json:"id" gorm:"primarykey"`
	WithdrawalID      uint `gorm:"index:wd_valid_unique,unique"`
	OrderAdjustmentID uint `gorm:"index:wd_valid_unique,unique"`
}

type WdOrderNotFound struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	OrderRefID string    `json:"order_ref_id"`
	WdID       uint      `json:"wd_id"`
	Amount     float64   `json:"amount"`
	At         time.Time `json:"at"`

	Wd *Withdrawal `json:"withdrawal"`
}
