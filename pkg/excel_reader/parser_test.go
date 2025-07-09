package excel_reader_test

import (
	"testing"

	"github.com/pdcgo/shared/pkg/excel_reader"
	"github.com/stretchr/testify/assert"
)

type MockEnum string

type MockData struct {
	Dc     int     `xls:"0"`
	Flo    string  `xls:"1"`
	Ss     float32 `xls:"2"`
	DDr    int
	Status MockEnum `xls:"3"`
}

func TestUnmarshal(t *testing.T) {
	t.Run("testing normal", func(t *testing.T) {
		item := MockData{}
		dataraw := []string{"123", "assadasdasd", "12.7", "batal"}
		err := excel_reader.UnmarshalRow(&item, dataraw, excel_reader.MetaIndex{})
		assert.Nil(t, err)
		assert.NotEmpty(t, item)
		assert.NotEmpty(t, item.Status)

		// data, _ := json.MarshalIndent(item, "", "  ")
		// t.Error(string(data))
	})

	t.Run("jika bukan pointer", func(t *testing.T) {
		item := MockData{}
		dataraw := []string{"123", "assadasdasd", "12.7", "batal"}
		err := excel_reader.UnmarshalRow(item, dataraw, excel_reader.MetaIndex{})
		assert.NotNil(t, err)
	})
}

func TestMarshal(t *testing.T) {
	t.Run("test case marshal normal", func(t *testing.T) {
		item := MockData{
			Dc:     12,
			Flo:    "asdasd",
			Ss:     12.12,
			DDr:    1,
			Status: "efhaasde3qwgbatal",
		}

		hasil, err := excel_reader.MarshalRow(&item)
		assert.Nil(t, err)
		assert.NotEmpty(t, hasil)

		// t.Error(hasil)
	})
}

func TestExcelRef(t *testing.T) {
	ref := excel_reader.ExcelRef("AA23")
	c, d := ref.GetCol()
	assert.Equal(t, 27, c)
	assert.Equal(t, "AA", d)

	row, err := ref.GetRow()

	assert.Nil(t, err)
	assert.Equal(t, 23, row)
}
