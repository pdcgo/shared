package db_models

import (
	"time"

	"github.com/pdcgo/shared/interfaces/identity_iface"
)

type AssetType string

const (
	UnspentStock  AssetType = "unspent_stock"
	MpHoldingFund AssetType = "mp_holding_fund"
	ABankAccount  AssetType = "bank_account"
)

type Asset struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	TeamID      uint      `json:"team_id"`
	AssetType   AssetType `json:"asset_type"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	LastUpdated time.Time `json:"last_updated"`
	CreatedAt   time.Time `json:"created_at"`

	Team *Team `json:"team"`
}

type AssetHistoryType string

const (
	AssetTransfer   AssetHistoryType = "tf"
	AssetAdjustment AssetHistoryType = "adj"
	AssetFund       AssetHistoryType = "fund"
	MpFund          AssetHistoryType = "mp_fund"
)

type AssetHistory struct {
	ID          uint `json:"id" gorm:"primarykey"`
	CreatedByID uint `json:"created_by_id"`
	FromAssetID uint `json:"from_asset_id"`
	ToAssetID   uint `json:"to_asset_id"`

	Type    AssetHistoryType         `json:"type"`
	At      time.Time                `json:"at" gorm:"index"`
	Amount  float64                  `json:"amount"`
	From    identity_iface.AgentType `json:"from"`
	IsValid bool                     `json:"is_valid" gorm:"index"`

	CreatedBy *User  `json:"created_by"`
	FromAsset *Asset `json:"from_asset" gorm:"foreignkey:FromAssetID"`
	ToAsset   *Asset `json:"to_asset" gorm:"foreignkey:ToAssetID"`
}

// GetEntityID implements authorization.Entity.
func (a *AssetHistory) GetEntityID() string {
	return "asset_history"
}
