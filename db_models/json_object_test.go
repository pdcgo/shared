package db_models_test

import (
	"encoding/json"
	"testing"

	"github.com/pdcgo/shared/db_models"
	"github.com/stretchr/testify/assert"
)

type MockD struct {
	C int `json:"c"`
}

func TestJsonDatabaseStruct(t *testing.T) {
	t.Run("testing marshal", func(t *testing.T) {
		mockd := db_models.NewJSONType(&MockD{
			C: 10,
		})

		data, err := json.Marshal(mockd)
		assert.Nil(t, err)
		assert.Contains(t, string(data), `"c":10`)

		t.Run("testing unmarshalling", func(t *testing.T) {
			mockd := db_models.NewJSONType(&MockD{})
			err := json.Unmarshal(data, &mockd)
			assert.Nil(t, err)

			assert.Equal(t, 10, mockd.DataO.C)
		})
	})
}
