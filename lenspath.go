package lenspath

import (
	"fmt"
	"reflect"
)

type Lens = string

type Lenspath struct {
	lens []Lens
	//assumeNil bool // if lenspath cannot be resolved, assume nil. Only for "get" operations
}

func Create(lens []Lens) (*Lenspath, error) {
	if len(lens) == 0 {
		return nil, &EmptyLensPathErr{}
	}
	return &Lenspath{lens: lens}, nil
}

func (lp *Lenspath) Get(value any) (any, error) {
	return lp.get(value, 0)
}

func (lp *Lenspath) get(value any, view int) (any, error) {
	switch {
	case view == lp.len():
		return value, nil

	case value == nil:
		return nil, &InvalidLensPathErr{}

	default:
		return lp.redirect(value, view)
	}
}

func (lp *Lenspath) redirect(value any, view int) (any, error) {
	if v, ok := value.(map[string]any); ok {
		return lp.getFromMap(v, view)
	} else if reflect.TypeOf(value).Kind() == reflect.Slice || reflect.TypeOf(value).Kind() == reflect.Array {
		if lp.path(view) == "*" {
			arr := reflect.ValueOf(value)
			slice := make([]interface{}, arr.Len())
			for j := 0; j < arr.Len(); j++ {
				slice[j] = arr.Index(j).Interface()
			}
			return lp.getAllFromList(slice, view+1)
		} else {
			return nil, NewInvalidLensPathErr(view, ArrayExpectedErr)
		}
	} else if reflect.TypeOf(value).Kind() == reflect.Struct {
		nestv := reflect.ValueOf(value).FieldByName(lp.path(view))
		if nestv.IsZero() {
			return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
		}

		return lp.get(nestv.Interface(), view+1)
	} else if reflect.TypeOf(value).Kind() == reflect.Ptr {
		return lp.get(reflect.ValueOf(value).Elem().Interface(), view)
	} else {
		return nil, fmt.Errorf("TODO: unhandled case: %T", value)
	}
}

func (lp *Lenspath) getAllFromList(value []any, view int) ([]any, error) {
	arr := make([]any, 0, len(value))
	for _, v := range value {
		if v, err := lp.get(v, view); err == nil {
			arr = append(arr, v)
		} else {
			return nil, err
		}
	}

	return arr, nil
}

func (lp *Lenspath) getFromMap(value map[string]any, view int) (any, error) {
	if mpv, ok := value[string(lp.lens[view])]; !ok {
		return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
	} else if val, err := lp.get(mpv, view+1); err != nil {
		return nil, err
	} else {
		return val, nil
	}
}

func (lp *Lenspath) len() int {
	return len(lp.lens)
}

func (lp *Lenspath) path(view int) string {
	return string(lp.lens[view])
}
