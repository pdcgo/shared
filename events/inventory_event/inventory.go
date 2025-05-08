package inventory_event

type Action string

const (
	AcceptStock Action = "accept_stock"
)

type TxEvent struct {
	Action Action `json:"action"`
	TxID   uint   `json:"inventory_id"`
}

// EventPath implements streampipe.Event.
func (t *TxEvent) EventPath() string {
	return "inventory/tx"
}
