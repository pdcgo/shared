package wallet_iface

import (
	"time"

	"github.com/pdcgo/shared/db_models"
)

type WalletStatus string

const (
	WalletActive    WalletStatus = "active"
	WalletSuspended WalletStatus = "suspended"
)

type BalanceType string

const (
	StockBalance BalanceType = "stock_balance"
	OpsBalance   BalanceType = "ops_balance"
)

type Wallet struct {
	ID          uint         `json:"id" gorm:"primarykey"`
	TeamID      uint         `json:"team_id" gorm:"index:id_balance_type,unique"`
	BalanceType BalanceType  `json:"balance_type" gorm:"index:id_balance_type,unique"`
	Balance     float64      `json:"balance"`
	Status      WalletStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Team db_models.Team
}

type TransactionType string

const (
	TransactionCredit TransactionType = "credit"
	TransactionDebit  TransactionType = "debit"
)

type TxType string

const (
	TxTypeBrokenGood      TxType = "broken_good"
	TxTypeBrokenGoodFound TxType = "broken_good_found"
	TxTypeUnknown         TxType = "unknown"
)

type TxStatus string

const (
	TxCompleted TxStatus = "completed"
)

type WalletTransaction struct {
	ID          uint            `json:"id" gorm:"primarykey"`
	TxID        string          `json:"tx_id" gorm:"index"`
	WalletID    uint            `json:"wallet_id"`
	Type        TransactionType `json:"type"`
	TxType      TxType          `json:"tx_type"`
	Status      TxStatus        `json:"status"`
	ReferenceID string          `json:"reference_id"`

	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	Wallet Wallet `json:"wallet"`
}

type WalletTransactionLog struct {
	// 	id	UUID / INT	Primary key
	// transaction_id	UUID / INT	FK to Transactions
	// status	ENUM	"success", "failed", "reversed"
	// logged_at	TIMESTAMP	Time of audit log
	// notes	TEXT	Additional info
}
