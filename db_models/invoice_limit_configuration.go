package db_models

import "github.com/pdcgo/shared/interfaces/invoice_iface"

// Deprecated: sudah ada owe configuration yang baru
type InvoiceLimitConfiguration struct {
	ID        uint                    `gorm:"primarykey" json:"id"`
	LimitType invoice_iface.LimitType `json:"limit_type"`
	TeamID    int64                   `json:"team_id" gorm:"uniqueIndex:team_for_unique"`
	ForTeamID *int64                  `json:"for_team_id" gorm:"uniqueIndex:team_for_unique"`
	Threshold float64                 `json:"threshold"`
}
