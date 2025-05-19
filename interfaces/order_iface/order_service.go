package order_iface

import "github.com/pdcgo/shared/db_models"

type OrderTagMutation interface {
	Add(from db_models.RelationFrom, orderIDs []uint, tags []string) error
	Remove(from db_models.RelationFrom, orderIDs []uint, tags []string) error
	RemoveAllFrom(from db_models.RelationFrom, orderIDs []uint) error
}
