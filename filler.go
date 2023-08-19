package anyrest

import (
	"errors"
	"fmt"
	"reflect"
)

type Filler struct {
}

func FillStruct(modelRef *interface{}, fields []any, values []any) (errs error) {
	model := *modelRef

	for i, name := range fields {
		err := error(nil)

		if reflect.TypeOf(name).Kind() == reflect.Slice {
			row := name.([]any)
			inners := row[1].([]any)

			fld := reflect.ValueOf(model).Elem().FieldByName(row[0].(string))
			switch fld.Type().Kind() {
			case reflect.Struct:
				field := fld.Addr().Interface()
				err = FillStruct(&field, inners, values[i].([]any))
			case reflect.Slice:
				slType := fld.Type().Elem()
				sl := reflect.MakeSlice(reflect.SliceOf(slType), 0, 0)
				for _, record := range values[i].([]any) {
					item := reflect.New(slType).Interface()
					err = FillStruct(&item, inners, record.([]any))
					sl = reflect.Append(sl, reflect.ValueOf(item).Elem())
				}
				fld.Set(sl)
			}

		} else {
			err = SetField(model, name.(string), values[i])
		}

		if err != nil {
			if errs == nil {
				errs = err
			} else {
				errors.Join(errs, err)
			}
		}
	}
	return errs
}

func SetField(m interface{}, name string, value any) error {
	rfv := reflect.ValueOf(value)
	if !rfv.IsValid() {
		return nil
	}

	rm := reflect.ValueOf(m)
	if rm.Kind() != reflect.Ptr || rm.Elem().Kind() != reflect.Struct {
		return errors.New("Model must be pointer to struct")
	}

	// Dereference pointer
	rm = rm.Elem()

	fv := rm.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("Unknown field name: %s", name)
	}

	if !fv.CanSet() {
		return fmt.Errorf("cannot set field %s", name)
	}

	if !rfv.CanConvert(fv.Type()) {
		return fmt.Errorf("%s of type %s can't convert to type %s", name, rfv.Kind(), fv.Kind())
	}

	cv := reflect.ValueOf(value).Convert(fv.Type())

	switch fv.Kind() {
	case reflect.String:

		fv.SetString(cv.String())

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		if reflect.ValueOf(value).Convert(reflect.TypeOf(0)).Int() < 0 {
			return fmt.Errorf("%s less 0 and can't convert to type %s", name, fv.Kind())
		}
		fv.SetUint(cv.Uint())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:

		fv.SetInt(cv.Int())

	case reflect.Float32, reflect.Float64:

		fv.SetFloat(cv.Float())

	case reflect.Bool:

		fv.SetBool(cv.Bool())

	default:
		return fmt.Errorf("%s has not an available type %s", name, fv.Kind())

	}

	return nil
}
