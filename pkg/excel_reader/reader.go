package excel_reader

import (
	"archive/zip"
	"errors"
	"io"
	"regexp"
)

var ErrFileNotFound = errors.New("file in gzip not found")

type CustomReadCloser zip.Reader

func (reader *CustomReadCloser) Blob(name string) ([]string, error) {
	ssrc := []string{}
	var err error

	reg, err := regexp.Compile(name)
	if err != nil {
		return ssrc, err
	}

	for _, v := range reader.File {
		if !reg.Match([]byte(v.Name)) {
			continue
		}

		ssrc = append(ssrc, v.Name)

	}

	return ssrc, nil
}

func (reader *CustomReadCloser) ReadFile(fname string) (io.ReadCloser, error) {
	var ssrc io.ReadCloser
	var err error
	for _, v := range reader.File {
		if v.Name != fname {
			continue
		}
		ssrc, err = v.Open()
		if err != nil {
			return ssrc, err
		}

	}

	if ssrc == nil {
		err = ErrFileNotFound
	}

	return ssrc, err
}

type ExcelReader struct {
	file *CustomReadCloser
}

func NewExcelReader(file *zip.Reader) *ExcelReader {
	reader := ExcelReader{
		file: (*CustomReadCloser)(file),
	}
	// shared, err := reader.loadSharedString()
	// if err != nil {
	// 	return &reader, err
	// }
	// reader.shared = shared

	return &reader
}

func (read *ExcelReader) GetWorkbook() (*Workbook, error) {
	return NewWorkbook(read.file)
}
