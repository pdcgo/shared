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

type Wallet struct {
	ID      uint         `json:"id" gorm:"primarykey"`
	TeamID  uint         `json:"team_id" gorm:"unique"`
	Balance float64      `json:"balance"`
	Status  WalletStatus `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Team db_models.Team
}

type WalletTransaction struct {
	// 	id	UUID / INT	Primary key
	// wallet_id	UUID / INT	FK to Wallets
	// type	ENUM	"credit", "debit"
	// amount	DECIMAL(12,2)	Always positive
	// description	TEXT	Optional note or reference
	// reference_id	VARCHAR	External ref (e.g., order ID, payment ID)
	// created_at	TIMESTAMP	Time of transaction

}

type WalletTransactionLog struct {
	// 	id	UUID / INT	Primary key
	// transaction_id	UUID / INT	FK to Transactions
	// status	ENUM	"success", "failed", "reversed"
	// logged_at	TIMESTAMP	Time of audit log
	// notes	TEXT	Additional info
}
