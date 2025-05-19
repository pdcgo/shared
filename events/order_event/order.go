package order_event

import "github.com/pdcgo/shared/db_models"

type Action string

const (
	OrderCreated      Action = "created"
	OrderChangeStatus Action = "change_status"
)

type OrderEvent struct {
	Action  Action              `json:"action"`
	OrderID uint                `json:"inventory_id"`
	Status  db_models.OrdStatus `json:"status"`
}

// EventPath implements streampipe.Event.
func (t *OrderEvent) EventPath() string {
	return "order"
}
