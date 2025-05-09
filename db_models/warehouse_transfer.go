package db_models

import (
	"time"
)

type WarehouseTransfer struct {
	ID              uint `json:"id" gorm:"primarykey"`
	InboundTxID     uint `json:"inbound_tx_id"`
	OutboundTxID    uint `json:"outbound_tx_id"`
	FromWarehouseID uint `json:"from_warehouse_id"`
	ToWarehouseID   uint `json:"to_warehouse_id"`
	TeamID          uint `json:"team_id"`

	Status InvTxStatus `json:"status"`

	InboundTx     *InvTransaction `json:"inbound"`
	OutboundTx    *InvTransaction `json:"outbound"`
	FromWarehouse *Warehouse      `json:"from_warehouse"`
	ToWarehouse   *Warehouse      `json:"to_warehouse"`
	Team          *Team           `json:"team"`
	CreatedAt     time.Time       `json:"created_at"`
}

// GetEntityID implements authorization.Entity.
func (w *WarehouseTransfer) GetEntityID() string {
	return "warehouse_transfer"
}
