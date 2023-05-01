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

type EmptyLensPathErr struct{}

func (e *EmptyLensPathErr) Error() string {
	return "lenspath: must have at least one lens"
}

type InvalidLensPathErr struct {
	index int
}

func (e *InvalidLensPathErr) Error() string {
	return fmt.Sprintf("lenspath: could not navigate further, end of structure reached at %d", e.index)
}

func Create(lens []Lens) (*Lenspath, error) {
	if len(lens) == 0 {
		return nil, &EmptyLensPathErr{}
	}
	return &Lenspath{lens: lens}, nil
}

func (lp *Lenspath) Get(value interface{}) (interface{}, error) {
	return lp.get(value, 0)
}

func (lp *Lenspath) get(value interface{}, view int) (interface{}, error) {
	switch {
	case view == lp.len():
		return value, nil

	case value == nil:
		return nil, &InvalidLensPathErr{}

	default:
		return lp.redirect(value, view)
	}
}

func (lp *Lenspath) redirect(value interface{}, view int) (interface{}, error) {
	if v, ok := value.(map[string]interface{}); ok {
		return lp.getFromMap(v, view)
	} else if _, ok := value.([]interface{}); ok {
		return nil, fmt.Errorf("TODO: unhandled array")
	} else if reflect.TypeOf(value).Kind() == reflect.Struct {
		nestv := reflect.ValueOf(value).FieldByName(lp.path(view))
		if nestv == (reflect.Value{}) {
			return nil, &InvalidLensPathErr{index: view}
		}

		return lp.get(nestv.Interface(), view+1)
	} else if reflect.TypeOf(value).Kind() == reflect.Ptr {
		return lp.get(reflect.ValueOf(value).Elem().Interface(), view)
	} else {
		return nil, fmt.Errorf("TODO: unhandled case: %T", value)
	}
}

func (lp *Lenspath) getAllFromList(value []interface{}, view int) ([]interface{}, error) {
	arr := make([]interface{}, len(value))
	for _, v := range value {
		if v, err := lp.get(v, view); err == nil {
			arr = append(arr, v)
		} else {
			return nil, err
		}
	}

	return arr, nil
}

func (lp *Lenspath) getFromMap(value map[string]interface{}, view int) (interface{}, error) {
	return lp.get(value[string(lp.lens[view])], view+1)
}

func (lp *Lenspath) len() int {
	return len(lp.lens)
}

func (lp *Lenspath) path(view int) string {
	return string(lp.lens[view])
}
