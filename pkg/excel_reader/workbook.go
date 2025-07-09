package excel_reader

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

// ---------------------- type workbook ------------------------------

type Workbook struct {
	Sheets map[string]*Sheet
	reader *CustomReadCloser
}

func NewWorkbook(reader *CustomReadCloser) (*Workbook, error) {
	file, err := reader.ReadFile("xl/workbook.xml")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	dec := xml.NewDecoder(file)
	workbook := Workbook{
		Sheets: map[string]*Sheet{},
		reader: reader,
	}

	// load shared string
	shared, err := NewSharedString(reader)
	if err != nil {
		return nil, err
	}

	var sheets []string
	sheets, err = reader.Blob(`xl\/worksheets/\/*`)
	if err != nil {
		return nil, err
	}
	c := 0
	for {
		token, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Inspect the type of the token just read.
		switch se := token.(type) {
		case xml.StartElement:

			ele := se.Name.Local
			if ele != "sheet" {
				continue
			}

			name := ""
			id := ""

			for _, d := range se.Attr {

				if name != "" && id != "" {
					continue
				}

				switch d.Name.Local {
				case "name":
					name = d.Value
				case "sheetId":
					id = d.Value
				case "rId":
					id = d.Value
				case "id":
					val := strings.ReplaceAll(d.Value, "rId", "")
					id = val
				}
			}

			if name == "" || id == "" {
				continue
			}
			sheet := &Sheet{
				Fname:     sheets[c],
				SharedStr: shared,
				reader:    reader,
			}
			workbook.Sheets[name] = sheet
			c += 1

		default:
			continue
		}
	}

	return &workbook, nil
}

func (work *Workbook) GetSheet(name string) (*Sheet, error) {
	if work.Sheets[name] == nil {
		return nil, fmt.Errorf("sheet not found for name %s", name)
	}

	return work.Sheets[name], nil
}
