package db_connect

import (
	"github.com/pdcgo/schema/services/common/v1"
	"gorm.io/gorm"
)

func SetPaginationQuery(
	db *gorm.DB,
	queryBuilder func() (*gorm.DB, error),
	pagination *common.PageFilter,
) (
	*gorm.DB,
	*common.PageInfo,
	error,
) {
	offset := (pagination.Page - 1) * pagination.Limit
	query, err := queryBuilder()
	if err != nil {
		return nil, nil, err
	}

	paginated := query.
		Offset(int(offset)).
		Limit(int(pagination.Limit))

	var totalItems int64

	qp, err := queryBuilder()
	if err != nil {
		return nil, nil, err
	}

	err = db.
		Table("(?) as d", qp).
		Select("count(1)").
		Find(&totalItems).
		Error

	if err != nil {
		return nil, nil, err
	}

	info := &common.PageInfo{
		CurrentPage: pagination.Page,
		TotalItems:  totalItems,
	}

	if totalItems <= pagination.Limit {
		info.TotalPage = 1
	} else {
		info.TotalPage = totalItems / pagination.Limit
		if totalItems%pagination.Limit != 0 {
			info.TotalPage++
		}
	}

	return paginated, info, nil
}
