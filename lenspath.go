package lenspath

import (
	"fmt"
	"reflect"
)

type Lens = string

type Lenspath struct {
	lens      []Lens
	assumeNil bool // if lenspath cannot be resolved, assume nil. If false, return error on unresolved lenspath while traversing structures
}

func Create(lens []Lens) (*Lenspath, error) {
	if len(lens) == 0 {
		return nil, &EmptyLensPathErr{}
	}
	return &Lenspath{lens: lens, assumeNil: true}, nil
}

func (lp *Lenspath) Get(data any) (any, error) {
	return lp.get(data, 0)
}

func (lp *Lenspath) Set(data any, value any) (any, error) {
	return lp.set(data, value, 0)
}

func (lp *Lenspath) get(data any, view int) (any, error) {
	if view == lp.len() {
		return data, nil
	} else if data == nil {
		if lp.assumeNil {
			return nil, nil
		} else {
			return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
		}
	}

	kind := reflect.TypeOf(data).Kind()

	switch kind {
	case reflect.Map:
		return lp.getFromMap(data, view)

	case reflect.Slice, reflect.Array:
		if lp.path(view) == "*" {

			// return []any if the array is not homogeneous (some lens gets return nil for example
			// or the map entries have different types for same keys)
			// else if array is homogeneous, return []<type> (e.g. []string)

			arr := reflect.ValueOf(data)
			any_slice := make([]any, 0, arr.Len())
			consistent_type := true
			var prev_type reflect.Type

			for j := 0; j < arr.Len(); j++ {
				if value, err := lp.get(arr.Index(j).Interface(), view+1); err == nil {
					value_type := reflect.TypeOf(value)
					any_slice = append(any_slice, value)
					if j > 0 {
						consistent_type = consistent_type && value_type == prev_type
					}
					prev_type = value_type
				} else {
					return nil, err
				}
			}

			if consistent_type {
				slice := reflect.MakeSlice(reflect.SliceOf(prev_type), 0, arr.Len())
				for _, v := range any_slice {
					slice = reflect.Append(slice, reflect.ValueOf(v))
				}

				return slice.Interface(), nil
			}

			return any_slice, nil
		} else {
			return nil, NewInvalidLensPathErr(view, ArrayExpectedErr)
		}

	case reflect.Struct:
		nestv := reflect.ValueOf(data).FieldByName(lp.path(view))
		if !nestv.IsValid() || nestv.IsZero() {
			if lp.assumeNil {
				return nil, nil
			} else {
				return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
			}
		}

		return lp.get(nestv.Interface(), view+1)

	case reflect.Ptr:
		return lp.get(reflect.ValueOf(data).Elem().Interface(), view)

	default:
		return nil, fmt.Errorf("unhandled case: %T", data)
	}
}

func (lp *Lenspath) set(data any, value any, view int) (any, error) {
	if view == lp.len() {
		return value, nil
	}

	kind := reflect.TypeOf(data).Kind()

	switch kind {
	case reflect.Map:
		return lp.setFromMap(data, value, view)

	case reflect.Slice, reflect.Array:
		if lp.path(view) == "*" {
			arr := reflect.ValueOf(data)
			fmt.Println(arr)
			slice := reflect.MakeSlice(arr.Type(), 0, arr.Len())

			// check if value is a slice or array; the length should then match
			// each value in the array is set to the corresponding value in the data slice
			if reflect.TypeOf(value).Kind() != reflect.Slice && reflect.TypeOf(value).Kind() != reflect.Array {
				return nil, ArrayParamExpectedErr
			}
			value_arr := reflect.ValueOf(value)

			if arr.Len() != value_arr.Len() {
				return nil, ParamSizeMismatchErr
			}

			for j := 0; j < arr.Len(); j++ {
				if v, err := lp.set(arr.Index(j).Interface(), value_arr.Index(j).Interface(), view+1); err == nil {
					slice = reflect.Append(slice, reflect.ValueOf(v))
				} else {
					return nil, err
				}
			}
			return slice.Interface(), nil
		} else {
			return nil, NewInvalidLensPathErr(view, ArrayExpectedErr)
		}

	case reflect.Struct:
		field := reflect.ValueOf(data).FieldByName(lp.path(view))
		if field.IsZero() {
			if lp.assumeNil {
				return nil, nil
			} else {
				return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
			}
		}

		if field.CanSet() {
			if val, err := lp.set(field.Interface(), value, view+1); err != nil {
				return nil, err
			} else {
				field.Set(reflect.ValueOf(val))
			}
		}

		return data, nil

	case reflect.Ptr:
		return lp.set(reflect.ValueOf(data).Elem().Interface(), value, view)

	default:
		return nil, fmt.Errorf("unhandled case: %T", data)
	}
}

func (lp *Lenspath) setFromMap(data any, value any, view int) (any, error) {
	key := reflect.ValueOf((lp.lens[view]))
	keyv := reflect.ValueOf(data).MapIndex(key)

	tosetv := value
	if !keyv.IsValid() || keyv.IsZero() {
		if view < lp.len()-1 {
			return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
		}
	} else if val, err := lp.set(keyv.Interface(), value, view+1); err != nil {
		return nil, err
	} else {
		tosetv = val
	}

	reflect.ValueOf(data).SetMapIndex(key, reflect.ValueOf(tosetv))
	return data, nil
}

func (lp *Lenspath) getFromMap(value any, view int) (any, error) {
	key := reflect.ValueOf((lp.lens[view]))
	keyv := reflect.ValueOf(value).MapIndex(key)

	if !keyv.IsValid() || keyv.IsZero() {
		if lp.assumeNil {
			return nil, nil
		} else {
			return nil, NewInvalidLensPathErr(view, LensPathStoppedErr)
		}
	} else {
		return lp.get(keyv.Interface(), view+1)
	}
}

func (lp *Lenspath) len() int {
	return len(lp.lens)
}

func (lp *Lenspath) path(view int) string {
	return string(lp.lens[view])
}
