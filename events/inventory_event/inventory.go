package inventory_event

import "github.com/pdcgo/shared/db_models"

type Action string

const (
	RestockAccept   Action = "restock_accept"
	RestockCreated  Action = "restock_created"
	RestockCancel   Action = "restock_cancel"
	InvChangeStatus Action = "change_status"
)

type TxEvent struct {
	Action Action                `json:"action"`
	TxID   uint                  `json:"inventory_id"`
	Type   db_models.InvTxType   `json:"type"`
	Status db_models.InvTxStatus `json:"status"`
}

// EventPath implements streampipe.Event.
func (t *TxEvent) EventPath() string {
	return "inventory/tx"
}
