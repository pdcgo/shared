package excel_reader

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ExcelRef string

var collist string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (ref ExcelRef) GetCol() (int, string) {

	coll := ""
	for _, ch := range ref {
	CC:
		for _, col := range collist {
			if ch == col {
				coll += string(ch)
				break CC
			}
		}
	}

	hasil := 0
	c := len(coll)
	po := 0
	for {
		c -= 1
		if c < 0 {
			break
		}

		idx := strings.Index(collist, string(coll[c])) + 1

		hasil += idx * int(math.Pow(26, float64(po)))
		po += 1

	}

	return hasil, coll
}

func (ref ExcelRef) GetRow() (int, error) {
	coll := ""
	for _, ch := range ref {
		found := false
	CC:
		for _, col := range collist {
			if ch == col {
				found = true
				break CC
			}
		}

		if !found {
			coll += string(ch)
		}
	}

	return strconv.Atoi(coll)
}

type MetaIndex []string

func (m MetaIndex) GetIndex(key string) (int, error) {
	for i, s := range m {
		if s == key {
			return i, nil
		}
	}

	return 0, fmt.Errorf("failed find index %s", key)
}

func UnmarshalRow(val interface{}, row []string, metaHeaders MetaIndex) error {
	var err error
	if reflect.ValueOf(val).Kind() != reflect.Ptr {
		return errors.New("marshalling row must pointer supplied")
	}

	v := reflect.ValueOf(val).Elem()
	tipe := v.Type()

	rowlen := len(row)
	useMeta := len(metaHeaders) != 0
Parent:
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tip := tipe.Field(i)
		// value := row[i]

		// getting idx

		var idxMethod []int8

		if useMeta {
			idxMethod = []int8{1, 2}
		} else {
			idxMethod = []int8{2}
		}

		var idx int
	LMethod:
		for _, acint := range idxMethod {
			switch acint {
			case 1: // dengan key header
				headstr := tip.Tag.Get("xlsheader")
				if headstr == "" {
					continue LMethod
				}
				idx, err = metaHeaders.GetIndex(headstr)
				if err != nil {
					return err
				}
				break LMethod
			case 2: // dengan key xls
				idxs := tip.Tag.Get("xls")
				if idxs == "" {
					continue Parent
				}

				idx, err = strconv.Atoi(idxs)
				if err != nil {
					return err
				}
				break LMethod
			}
		}

		if idx >= rowlen {
			return errors.New("excelrow cant parse because row length supplied not length enough")
		}
		data := row[idx]

		switch field.Kind() {
		case reflect.String:
			field.SetString(data)
		case reflect.Float32:
			var pdata float64 = 0
			pdata, err = strconv.ParseFloat(data, 64)
			field.SetFloat(pdata)
		case reflect.Float64:
			var pdata float64 = 0
			pdata, err = strconv.ParseFloat(data, 64)
			field.SetFloat(pdata)
		case reflect.Bool:
			data = strings.ToLower(data)
			pdata := false
			if data == "1" ||
				data == "true" {
				pdata = true
			}

			field.SetBool(pdata)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var pdata int = 0
			pdata, err = strconv.Atoi(data)
			field.SetInt(int64(pdata))

		default:
			switch field.Interface().(type) {
			case time.Time:
				dateformat := tip.Tag.Get("xlsdate")
				fallbackdate := tip.Tag.Get("fallback_xlsdate")
				addhour := tip.Tag.Get("addhour")
				if dateformat == "" {
					err = fmt.Errorf("%s not have format date", tip.Name)
					break
				}

				// localization := time.Local
				// if timeutc == "true" {
				// 	localization = time.UTC
				// }

				var t time.Time
				t, err = time.ParseInLocation(dateformat, data, time.Local)
				if err != nil && fallbackdate != "" {
					t, err = time.ParseInLocation(fallbackdate, data, time.Local)
				}

				if addhour != "" {
					t = t.Add(time.Hour * 23)
				}

				field.Set(reflect.ValueOf(t))
			}

		}

		if err != nil {
			return err
		}
	}

	return err
}

func RowNeed(val interface{}) (int, []int, error) {
	var err error
	maxArrayLen := 0
	indexNeed := []int{}

	v := reflect.ValueOf(val).Elem()
	tipe := v.Type()

	// checking max num length
	for i := 0; i < v.NumField(); i++ {
		tip := tipe.Field(i)
		idxs := tip.Tag.Get("xls")
		if idxs == "" {
			continue
		}
		var idx int
		idx, err = strconv.Atoi(idxs)
		if err != nil {
			return maxArrayLen, indexNeed, err
		}
		indexNeed = append(indexNeed, idx)
		if idx > maxArrayLen {
			maxArrayLen = idx
		}
	}

	return maxArrayLen, indexNeed, nil
}

func MarshalRow(val interface{}) ([]string, error) {
	var err error
	if reflect.ValueOf(val).Kind() != reflect.Ptr {
		return []string{}, errors.New("marshalling row must pointer supplied")
	}

	v := reflect.ValueOf(val).Elem()
	tipe := v.Type()

	maxArrayLen, _, err := RowNeed(val)
	if err != nil {
		return []string{}, err
	}

	if maxArrayLen == 0 {
		return []string{}, nil
	}

	hasil := make([]string, maxArrayLen+1)
	// parsing hasilnya
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tip := tipe.Field(i)
		idxs := tip.Tag.Get("xls")
		if idxs == "" {
			continue
		}
		var idx int
		idx, err = strconv.Atoi(idxs)
		if err != nil {
			return []string{}, err
		}
		if idx > maxArrayLen {
			maxArrayLen = idx
		}

		switch field.Kind() {
		case reflect.String:
			hasil[idx] = field.String()
		case reflect.Float32, reflect.Float64:
			hasil[idx] = fmt.Sprintf("%.5f", field.Float())

		case reflect.Bool:
			if field.Bool() {
				hasil[idx] = "1"
			} else {
				hasil[idx] = "0"
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			hasil[idx] = fmt.Sprintf("%d", field.Int())
		}
	}

	return hasil, err
}
