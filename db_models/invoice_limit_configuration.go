package db_models

import "github.com/pdcgo/shared/interfaces/invoice_iface"

type InvoiceLimitConfiguration struct {
	ID        uint                    `gorm:"primarykey" json:"id"`
	LimitType invoice_iface.LimitType `json:"limit_type"`
	TeamID    int64                   `json:"team_id"`
	ForTeamID *int64                  `json:"for_team_id"`
	Threshold float64                 `json:"threshold"`
}
