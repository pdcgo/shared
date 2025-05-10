package db_models_test

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/pdcgo/shared/db_models"
	"github.com/pdcgo/shared/pkg/debugtool"
	"github.com/stretchr/testify/assert"
)

func TestBasicModel(t *testing.T) {

	t.Run("testing normal", func(t *testing.T) {
		sku := db_models.Sku{
			VariantID:   3423,
			ProductID:   45,
			WarehouseID: 3,
			TeamID:      2,
		}

		hasil, err := sku.CalculateID()
		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)
	})

	t.Run("testing ada yang kosong", func(t *testing.T) {
		sku := db_models.Sku{
			VariantID:   3423,
			ProductID:   45,
			WarehouseID: 3,
		}

		hasil, err := sku.CalculateID()
		assert.NotNil(t, err)
		assert.Empty(t, hasil)
	})

	t.Run("testing basic model tanggal", func(t *testing.T) {
		inv := db_models.InvTransaction{
			Created: time.Now(),
		}

		t.Run("testing arrived null bukan kosong", func(t *testing.T) {
			data, err := json.Marshal(inv)
			assert.Nil(t, err)

			content := string(data)
			assert.Contains(t, content, `"arrived":null`)

		})

	})

}

func TestSortingHistories(t *testing.T) {
	// map 1 -- 13
	// map 2 -- 10

	histories := []*db_models.InvertoryHistory{
		{
			RackID: 1,
			Count:  -9,
		},
		{
			RackID: 2,
			Count:  -6,
		},
		{
			RackID: 1,
			Count:  -3,
		},
		{
			RackID: 1,
			Count:  -1,
		},
		{
			RackID: 2,
			Count:  -4,
		},
	}

	histories = db_models.InvHistorySort(histories)

	expected := []*db_models.InvertoryHistory{
		{
			RackID: 2,
			Count:  -4,
		},
		{
			RackID: 2,
			Count:  -6,
		},
		{
			RackID: 1,
			Count:  -1,
		},
		{
			RackID: 1,
			Count:  -3,
		},
		{
			RackID: 1,
			Count:  -9,
		},
	}

	for i, dd := range expected {
		assert.Equal(t, dd.RackID, histories[i].RackID, i)
		assert.Equal(t, dd.Count, histories[i].Count, i)
	}
}

func TestPartialModel(t *testing.T) {
	t.Run("testing partial invitem", func(t *testing.T) {
		tx := db_models.InvTransaction{
			Items: []*db_models.InvTxItem{
				{
					ID:               1,
					InvTransactionID: 1,
					SkuID:            "11111111",
					Count:            5,
					Price:            100,
					Total:            500,
				},
				{
					ID:               1,
					InvTransactionID: 1,
					SkuID:            "11111111",
					Count:            1,
					Price:            200,
					Total:            200,
				},
				{
					ID:               1,
					InvTransactionID: 1,
					SkuID:            "11111112",
					Count:            4,
					Price:            100,
					Total:            400,
				},
				{
					ID:               1,
					InvTransactionID: 1,
					SkuID:            "11111113",
					Count:            3,
					Price:            200,
					Total:            600,
				},
				{
					ID:               1,
					InvTransactionID: 1,
					SkuID:            "11111114",
					Count:            6,
					Price:            200,
					Total:            600,
				},
				{
					ID:               1,
					InvTransactionID: 1,
					SkuID:            "11111114",
					Count:            6,
					Price:            200,
					Total:            600,
				},
			},
		}

		hasil, err := tx.Items.GetItemPartials([]*db_models.PartialItem{
			{
				ProductID:   1,
				VariationID: 1,
				Count:       6,
			},
			{
				ProductID:   1,
				VariationID: 2,
				Count:       3,
			},
			{
				ProductID:   1,
				VariationID: 4,
				Count:       7,
			},
		})

		assert.Nil(t, err)

		sort.Slice(hasil, func(i, j int) bool {
			return hasil[i].Count < hasil[j].Count
		})

		checks := []*db_models.InvTxItem{
			{
				ID:               1,
				InvTransactionID: 1,
				SkuID:            "11111111",
				Count:            1,
				Price:            200,
				Total:            200,
			},
			{
				ID:               1,
				InvTransactionID: 1,
				SkuID:            "11111114",
				Count:            1,
				Price:            200,
				Total:            200,
			},
			{
				ID:               1,
				InvTransactionID: 1,
				SkuID:            "11111112",
				Count:            3,
				Price:            200,
				Total:            200,
			},

			{
				ID:               1,
				InvTransactionID: 1,
				SkuID:            "11111111",
				Count:            5,
				Price:            100,
				Total:            500,
			},
			{
				ID:               1,
				InvTransactionID: 1,
				SkuID:            "11111114",
				Count:            6,
				Price:            100,
				Total:            400,
			},
		}
		// dtas, _ := json.MarshalIndent(hasil, "", "  ")
		// t.Error(string(dtas))

		for i, item := range hasil {
			assert.Equal(t, checks[i].Count, item.Count)
		}

	})
}

func TestKaitanSortingHistoryGagalDIevent(t *testing.T) {

	histories := []*db_models.InvertoryHistory{

		{
			ID:          151070,
			RackID:      824,
			SkuID:       db_models.SkuID("2233263621626f"),
			WarehouseID: 38,
			TeamID:      54,
			UserID:      109,
			Count:       -1,
			Price:       36000,
			ExtPrice:    405,
			Created:     time.Now(),
		},
		{
			ID:          137563,
			RackID:      824,
			SkuID:       db_models.SkuID("2233263621626f"),
			WarehouseID: 38,
			TeamID:      54,
			UserID:      109,
			Count:       -55,
			Price:       41000,
			ExtPrice:    1062.5,
			Created:     time.Now(),
		},
		{
			ID:          141198,
			RackID:      824,
			SkuID:       db_models.SkuID("2233263621626f"),
			WarehouseID: 38,
			TeamID:      54,
			UserID:      109,
			Count:       -2,
			Price:       36000,
			ExtPrice:    405,
			Created:     time.Now(),
		},
	}

	// getting histories
	histories = db_models.InvHistorySort(histories)
	hist := histories[0]
	assert.Equal(t, 36405.00, hist.ExtPrice+hist.Price)
	debugtool.LogJson(histories)

}
