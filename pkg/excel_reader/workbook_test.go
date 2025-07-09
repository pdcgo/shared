package excel_reader_test

import (
	"archive/zip"
	"context"
	"log"
	"testing"

	"github.com/pdcgo/shared/pkg/excel_reader"
	"github.com/stretchr/testify/assert"
)

func TestWorkbookSheet(t *testing.T) {

	fnames := []string{
		"../../test/assets/wd.xlsx",
		"../../test/assets/tiktokincome.xlsx",
		"../../test/assets/order_cancel.xlsx",
		"../../test/assets/return.xls",
		"../../test/assets/my_balance_transaction_report.shopee.xlsx",
	}

	for _, fname := range fnames {
		f, err := zip.OpenReader(fname)
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		reader := excel_reader.NewExcelReader(&f.Reader)
		workbook, err := reader.GetWorkbook()
		assert.Nil(t, err)

		for shetname := range workbook.Sheets {
			sheet, err := workbook.GetSheet(shetname)
			assert.Nil(t, err)
			t.Log(fname, shetname)
			err = sheet.Iterate(context.Background(), func(row []string) error {
				return nil
			})

			assert.Nil(t, err)
		}
	}

}
