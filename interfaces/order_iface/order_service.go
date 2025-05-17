package order_iface

type RelationFrom string

const (
	RelationFromWarehouse RelationFrom = "warehouse"
	RelationFromUser      RelationFrom = "user"
	RelationFromTracking  RelationFrom = "tracking"
)

type OrderTagMutation interface {
	Add(from RelationFrom, orderIDs []uint, tags []string) error
	Remove(orderIDs []uint, tags []string) error
	RemoveAllFrom(from RelationFrom) error
}
