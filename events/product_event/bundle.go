package product_event

type Action string

const (
	ActionUpdate Action = "update"
	ActionCreate Action = "create"
	ActionDelete Action = "delete"
)

type BundleEvent struct {
	Action   Action `json:"action"`
	BundleID uint   `json:"bundle_id"`
}
