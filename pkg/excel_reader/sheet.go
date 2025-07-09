package excel_reader

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type StringItem struct {
	T string `xml:"t"`
}

type Col struct {
	TAttr string   `xml:"t,attr"`
	R     ExcelRef `xml:"r,attr"`
	Val   string   `xml:"v"`
	ValT  string   `xml:"is>t"`
}

type Row struct {
	Col []Col `xml:"c"`
}

type Sheet struct {
	Fname     string
	SharedStr *SharedString
	reader    *CustomReadCloser
}

func (sh *Sheet) IterWithInterface(ctx context.Context, item any, handle func(data []string, rowerr error) error) error {
	var err error

	fname := sh.Fname
	file, err := sh.reader.ReadFile(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := xml.NewDecoder(file)

	maxIdx, indexNeeds, err := RowNeed(item)
	if err != nil {
		return err
	}

	record := make([]*string, maxIdx+1)
	lastrow := 1

	for {
		token, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch se := token.(type) {
		case xml.StartElement:
			ele := se.Name.Local
			if ele != "c" {
				continue
			}

			var col Col
			secop := se.Copy()
			err = dec.DecodeElement(&col, &secop)
			if err != nil {
				return err
			}

			i, _ := col.R.GetCol()
			row, err := col.R.GetRow()
			if err != nil {
				return err
			}

			// checking row next
			if row < lastrow {
				return fmt.Errorf("missing data in row %d", row)
			}

			newrow := lastrow != row
			if newrow {
				record = make([]*string, maxIdx+1)
				lastrow = row
			}

			i = i - 1
			if i > maxIdx { // skip karena g butuh di struct
				continue
			}

			switch col.TAttr {
			case "inlineStr": // jika lazada dan nested xml di dalammya
				record[i] = &col.ValT
			case "n":
				record[i] = &col.Val
			case "s":
				c, err := strconv.Atoi(col.Val)
				if err != nil {
					return err
				}
				val := sh.SharedStr.Data[c]
				record[i] = &val
			default:
				record[i] = &col.Val
			}

			containnil := false
			for _, r := range indexNeeds {
				d := record[r]
				if d == nil {
					containnil = true
					break
				}
			}

			// log.Println(*record[i], row, containnil)

			// checking jika sudah tidak ada yang nil
			if !containnil {

				var rowerr error = nil
				if containnil {
					rowerr = fmt.Errorf("excelread lastrow %d invalid have empty value when read", lastrow)
				}

				// marshaling data
				datas := make([]string, maxIdx+1)
				for id, v := range record {
					if v == nil {
						continue
					}
					datas[id] = *v
				}

				// if datas[0] == "577212414310909526" && !containnil {
				// 	log.Println(datas, "asdasdasd", rowerr)
				// 	// panic("asdsd")
				// }
				// log.Println(sh.Fname, datas, rowerr)
				err = handle(datas, rowerr)
				if err != nil {
					return err
				}

			}

		}
	}

}

func (sh *Sheet) Iterate(ctx context.Context, handler func(row []string) error) error {

	fname := sh.Fname
	file, err := sh.reader.ReadFile(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	dec := xml.NewDecoder(file)

	for {
		token, err := dec.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		switch se := token.(type) {
		case xml.StartElement:
			ele := se.Name.Local
			if ele != "row" {
				continue
			}
			var row Row
			err = dec.DecodeElement(&row, &se)
			if err != nil {
				return err
			}

			// data, _ := json.MarshalIndent(row, "", "  ")
			// log.Println(string(data))

			rec := make([]string, len(row.Col))
			for i, v := range row.Col {
				if v.TAttr == "n" {
					rec[i] = v.Val
					continue
				}

				if v.TAttr == "s" {
					c, err := strconv.Atoi(v.Val)
					if err != nil {
						return err
					}

					rec[i] = sh.SharedStr.Data[c]
					continue
				}

				rec[i] = v.Val

			}
			select {
			case <-ctx.Done():
				return nil
			default:
				err = handler(rec)
				if err != nil {
					return err
				}
			}

			// rec := []string{}
			// for _, v := range row.C {
			// 	if v.T == "n" {
			// 		rec = append(rec, v.V)
			// 		continue
			// 	}
			// 	i, err := strconv.Atoi(v.V)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	rec = append(rec, r.sharedStrings[i])
			// }
			// return rec, nil
		default:
			continue
		}
	}

}
