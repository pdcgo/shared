package db_tools

import "gorm.io/gorm"

func FindInBatch(db *gorm.DB, dest interface{}, batchSize int, fc func(tx *gorm.DB, batch int) error) *gorm.DB {
	var (
		tx           = db.Session(&gorm.Session{})
		queryDB      = tx
		rowsAffected int64
		batch        int
	)

	for {
		result := queryDB.Limit(batchSize).Offset(batch * batchSize).Find(dest)
		rowsAffected += result.RowsAffected
		batch++

		if result.Error == nil && result.RowsAffected != 0 {
			tx.AddError(fc(result, batch))
		} else if result.Error != nil {
			tx.AddError(result.Error)
		}

		if tx.Error != nil || int(result.RowsAffected) < batchSize {
			break
		}
	}

	tx.RowsAffected = rowsAffected
	return tx
}
