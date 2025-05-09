package db_models_test

import (
	"testing"

	"github.com/pdcgo/shared/db_models"
	"github.com/stretchr/testify/assert"
)

func TestWarehouse(t *testing.T) {
	ware := db_models.Warehouse{
		UseFixedFee: true,
		FeeFix:      12000,
		FeePercent:  0.1,
		MaxFee:      50000,
	}

	t.Run("testing getting fee fix", func(t *testing.T) {
		price, err := ware.GetWarehouseFee(700000)
		assert.Nil(t, err)

		assert.Equal(t, float64(12000), price)

	})

	t.Run("testing percent fee", func(t *testing.T) {
		ware := db_models.Warehouse{
			UseFixedFee: false,
			FeeFix:      12000,
			FeePercent:  0.1,
			MaxFee:      15000,
		}

		price, err := ware.GetWarehouseFee(123456)
		assert.Nil(t, err)
		assert.Equal(t, float64(12400), price)

		t.Run("test getting fee not more than threshold", func(t *testing.T) {
			price, err := ware.GetWarehouseFee(200000) // if without threshold = 20100
			assert.Nil(t, err)
			assert.Equal(t, 15000, int(price))
		})
	})
}
