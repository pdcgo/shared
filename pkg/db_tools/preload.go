package db_tools

import (
	"gorm.io/gorm"
)

func Preload[D, K, P any, Field comparable](
	datas []D,
	keyfunc func(i int, data D) []K,
	query func(ids []K) (*gorm.DB, func(P) Field),
	setter func(i int, datamap map[Field]P),
) error {
	var err error

	ids := make([]K, 0)
	for i, data := range datas {
		key := keyfunc(i, data)
		ids = append(ids, key...)
	}

	preloads := make([]P, 0)
	q, fieldget := query(ids)

	err = q.
		Find(&preloads).
		Error

	if err != nil {
		return err
	}

	datamap := make(map[Field]P, 0)
	for _, d := range preloads {
		prel := d
		key := fieldget(prel)
		datamap[key] = prel
	}

	for i := range datas {
		setter(i, datamap)
	}

	return nil
}
