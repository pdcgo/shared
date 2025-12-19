package db_connect

import "gorm.io/gorm"

type NextFunc func(query *gorm.DB) (*gorm.DB, error)
type NextHandler func(db *gorm.DB, next NextFunc) NextFunc

func NewQueryChain(db *gorm.DB, chains ...NextHandler) (*gorm.DB, error) {

	var next NextFunc = func(query *gorm.DB) (*gorm.DB, error) {
		return query, nil
	}

	reverse(chains)
	for _, chain := range chains {
		next = chain(db, next)
	}

	return next(db)
}
func reverse[T any](s []T) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
