package excel_reader_test

import (
	"archive/zip"
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/pdcgo/shared/pkg/excel_reader"
	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	// fname := "../../test/assets/test.xlsx"
	// fname := "../../test/assets/order_cancel.xlsx"
	fname := "../../test/assets/return.xls"

	f, err := zip.OpenReader(fname)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := excel_reader.NewExcelReader(&f.Reader)
	workbook, err := reader.GetWorkbook()
	assert.Nil(t, err)

	sheet, err := workbook.GetSheet("Sheet1")
	// sheet, err := workbook.GetSheet("orders")

	assert.Nil(t, err)

	err = sheet.Iterate(context.Background(), func(row []string) error {
		data, _ := json.Marshal(row)
		log.Println(string(data))

		return nil
	})

	assert.Nil(t, err)
}

type TiktokWdItem struct {
	ExternalOrderID  string    `xls:"0"`
	Type             string    `xls:"1"`
	SettlementAmount float64   `xls:"5"`
	OrderSettledTime time.Time `xls:"3" xlsdate:"2006/01/02" timeutc:"true"`
}

func TestTiktokDataDesember(t *testing.T) {
	fname := "../../test/assets/testwd/pzen_tidak_complete.xlsx"

	f, err := zip.OpenReader(fname)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := excel_reader.NewExcelReader(&f.Reader)
	workbook, err := reader.GetWorkbook()
	assert.Nil(t, err)

	for key := range workbook.Sheets {
		sheet, err := workbook.GetSheet(key)
		assert.Nil(t, err)
		err = sheet.IterWithInterface(context.Background(), &TiktokWdItem{}, func(data []string, rowerr error) error {
			// data, _ := json.Marshal(row)
			// log.Println(data)

			return nil
		})

		assert.Nil(t, err)
	}

}
