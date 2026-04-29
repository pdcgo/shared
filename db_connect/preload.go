package db_connect

import (
	"errors"
	"reflect"

	"gorm.io/gorm"
)

type Foreign struct {
	Field              string
	RefKey             string
	Fkey               string
	ReferenceValueFunc func(i int) any
	QueryOpt           func(q *gorm.DB) *gorm.DB
}

// toSnakeCase converts a CamelCase string to snake_case.
// Consecutive uppercase acronyms are handled correctly:
//
//	UserID     → user_id
//	HTTPSServer → https_server
func toSnakeCase(s string) string {
	runes := []rune(s)
	var result []rune
	for i, r := range runes {
		if r >= 'A' && r <= 'Z' {
			if i > 0 {
				prevLower := runes[i-1] >= 'a' && runes[i-1] <= 'z'
				prevUpper := runes[i-1] >= 'A' && runes[i-1] <= 'Z'
				nextLower := i+1 < len(runes) && runes[i+1] >= 'a' && runes[i+1] <= 'z'
				if prevLower || (prevUpper && nextLower) {
					result = append(result, '_')
				}
			}
			result = append(result, r+('a'-'A'))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

func (f *Foreign) DbFkey() string { // to snakecase
	return toSnakeCase(f.Fkey)
}

func (f *Foreign) DbRefKey() string { // to snakecase
	return toSnakeCase(f.RefKey)
}

func OneToMany(
	db *gorm.DB,
	datas any,
	foreigner *Foreign,
) error {
	var err error
	mapper := map[any][]reflect.Value{}
	idmapper := map[any]bool{}

	listVal := reflect.ValueOf(datas)
	tipes := reflect.TypeOf(datas)

	switch tipes.Kind() {
	case reflect.Slice:
		// getting mapper key
		for i := 0; i < listVal.Len(); i++ {
			refId := foreigner.ReferenceValueFunc(i)
			item := listVal.Index(i)
			if item.Kind() == reflect.Pointer {
				item = reflect.Indirect(item)
			}
			mapper[refId] = append(mapper[refId], item)
			idmapper[refId] = true
		}
	default:
		return errors.New("data is not iterable")
	}

	idlist := []any{}
	for c := range idmapper {
		idlist = append(idlist, c)
	}

	elemzero := reflect.Zero(tipes.Elem()).Interface()
	itemsobj, err := getEmptyStruct(elemzero)
	if err != nil {
		return err
	}
	model, err := getFieldStruct(foreigner.Field, itemsobj)
	if err != nil {
		return err
	}

	refvalue := reflect.MakeSlice(reflect.SliceOf(model.Type()), 0, 0)

	// find index key in modellist
	var modellist = refvalue.Interface()

	if foreigner.QueryOpt == nil {
		return errors.New("query option is not defined")
	}

	querymod := foreigner.QueryOpt(db).
		Where(foreigner.DbFkey()+" IN (?)", idlist)

	err = querymod.
		Find(&modellist).
		Error

	if err != nil {
		return err
	}

	valmodelist := reflect.ValueOf(modellist)
	for i := 0; i < valmodelist.Len(); i++ {
		prevalitem := valmodelist.Index(i)

		if prevalitem.Kind() == reflect.Pointer {
			prevalitem = reflect.Indirect(prevalitem)
		}
		idkey := prevalitem.FieldByName(foreigner.RefKey).Interface()
		values := mapper[idkey]
		for _, item := range values {
			cpreload := item.FieldByName(foreigner.Field)

			switch cpreload.Kind() {
			case reflect.Slice:
				items := reflect.Append(cpreload, prevalitem.Addr())
				cpreload.Set(items)

			default:
				cpreload.Set(prevalitem.Addr())
			}
		}
	}

	return err

}

func getEmptyStruct(data interface{}) (any, error) {
	tipes := reflect.TypeOf(data)
	values := reflect.ValueOf(data)

	switch tipes.Kind() {
	case reflect.Pointer:
		value := reflect.Indirect(values)
		if values.IsNil() {
			val := reflect.Zero(tipes.Elem()).Interface()
			return val, nil
		}

		return value.Interface(), nil
	case reflect.Struct:
		return data, nil

	default:
		return nil, errors.New("kind not struct or pointer of struct")
	}

}

// getting model yang akan di preload
func getFieldStruct(field string, data any) (reflect.Value, error) {
	var err error

	tipes := reflect.TypeOf(data)
	values := reflect.ValueOf(data)

	for c := 0; c < tipes.NumField(); c++ {
		tipe := tipes.Field(c)
		val := values.Field(c)

		if tipe.Name == field {
			switch val.Kind() {
			case reflect.Slice:
				itemType := tipe.Type.Elem()
				val = reflect.Zero(itemType)
			}
			return val, nil
		}

	}
	err = errors.New("field " + field + " notfound")
	return reflect.Value{}, err
}
