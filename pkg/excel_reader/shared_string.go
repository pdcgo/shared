package excel_reader

import (
	"encoding/xml"
	"errors"
	"io"
)

type SharedString struct {
	Data []string
	// Data *xml.Decoder
}

func NewSharedString(reader *CustomReadCloser) (*SharedString, error) {
	sh := SharedString{
		Data: []string{},
	}

	file, err := reader.ReadFile("xl/sharedStrings.xml")
	if err != nil {
		if errors.Is(err, ErrFileNotFound) {
			return &sh, nil
		}
		return nil, err
	}

	defer file.Close()

	dec := xml.NewDecoder(file)

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
			if ele != "si" {
				continue
			}
			var si StringItem
			if err := dec.DecodeElement(&si, &se); err != nil {
				return nil, err
			}

			sh.Data = append(sh.Data, si.T)
		default:
			continue
		}
	}

	return &sh, nil
}
