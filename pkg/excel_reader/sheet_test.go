package excel_reader_test

import (
	"archive/zip"
	"context"
	"log"
	"testing"

	"github.com/pdcgo/shared/pkg/excel_reader"
	"github.com/stretchr/testify/assert"
)

type DumpItem struct {
	ExternalOrderID string `xls:"0"`
	Status          string `xls:"1"`
	SubStatus       string `xls:"2"`
	CancelType      string `xls:"3"`
	Slow            string `xls:"4"`
}

func TestSheet(t *testing.T) {
	fname := "../../test/assets/allpesanantitok.xlsx"
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

		t.Run("testing iterating yang baru", func(t *testing.T) {
			c := 1
			datar := []string{"", ""}
			err = sheet.IterWithInterface(context.Background(), &DumpItem{}, func(data []string, rowerr error) error {
				if data[0] == "577212414310909526" {
					datar = data

				}
				c += 1
				return nil
			})

			assert.Equal(t, "577212414310909526", datar[0])
			assert.Equal(t, "Canceled", datar[1])
			assert.Greater(t, c, 1)

			assert.Nil(t, err)
		})

	}
}

type TiktokReturnItem struct {
	ExternalOrderID string `xls:"1"`
	ReturnStatus    string `xls:"17"`
	ReturnSubStatus string `xls:"18"`
	TrackingID      string `xls:"16"`
	BuyerNote       string `xls:"24"`
}

func TestSheetDenganIncompleteRecord(t *testing.T) {
	fname := "../../test/assets/tiktokreturn.xlsx"
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

		t.Run("testing iterating yang baru", func(t *testing.T) {
			c := 0

			err = sheet.IterWithInterface(context.Background(), &TiktokReturnItem{}, func(data []string, rowerr error) error {
				if data[1] == "576745558605268140" {
					item := TiktokReturnItem{}
					err = excel_reader.UnmarshalRow(&item, data, excel_reader.MetaIndex{})
					assert.Nil(t, err)

					if err != nil {
						return err
					}
					// t.Error(data)
				}

				c += 1
				return nil
			})

			assert.Nil(t, err)
			assert.Equal(t, 96, c)
		})

	}
}
