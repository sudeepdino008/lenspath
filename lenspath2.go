package lenspath

import (
	"fmt"
	"reflect"
)

type Lens = string

type LeafCallback = func(any) any // called when a leaf (of the lenspath) is reached; return value is set to the leaf (in case of setter)

type Lenspath struct {
	lens         []Lens
	lastArrayPos int  // last position of array (*) lens in lenspath
	assumeNil    bool // if lenspath cannot be resolved, assume nil. If false, return error on unresolved lenspath while traversing structures
}

type TraversalDetails struct {
	callback LeafCallback
	settable bool
}

func Create(lens []Lens) (*Lenspath, error) {
	if len(lens) == 0 {
		return nil, &EmptyLensPathErr{}
	}
	lastArrPos := -1
	for i, lensv := range lens {
		if lensv == "*" {
			lastArrPos = i
		}
	}
	assumeNil := true // default to assume nil
	return &Lenspath{lens, lastArrPos, assumeNil}, nil
}

func (lp *Lenspath) Get(data any) (any, error) {
	if lp.isArrBased() {
		return nil, NewInvalidLensPathErr(-1, PathContainsArrErr)
	}
	var rdata any
	_, err := lp.recurse(data, 0, &TraversalDetails{callback: func(data any) any {
		rdata = data
		return nil
	}, settable: false})
	return rdata, err
}

func (lp *Lenspath) Getter(data any, callback LeafCallback) error {
	_, err := lp.recurse(data, 0, &TraversalDetails{callback: callback, settable: false})
	return err
}

func (lp *Lenspath) Set(data any, value any) error {
	if lp.isArrBased() {
		return NewInvalidLensPathErr(-1, PathContainsArrErr)
	}
	_, err := lp.recurse(data, 0, &TraversalDetails{callback: func(data any) any {
		return value
	}, settable: true})
	return err
}

func (lp *Lenspath) Setter(data any, callback LeafCallback) error {
	_, err := lp.recurse(data, 0, &TraversalDetails{callback: callback, settable: true})
	return err
}

func (lp *Lenspath) recurse(data any, view int, details *TraversalDetails) (any, error) {
	if view == lp.len() {
		return details.callback(data), nil
	} else if data == nil {
		return nil, nil
	}

	kind := reflect.TypeOf(data).Kind()

	switch kind {
	case reflect.Map:
		lp.traverseMap(data, view, details)

	case reflect.Slice, reflect.Array:
		if lp.path(view) == "*" {
			lp.traverseSlice(data, view, details)
		} else {
			return nil, NewInvalidLensPathErr(view, ArrayExpectedErr)
		}

	case reflect.Struct:
		nestv := reflect.ValueOf(data).FieldByName(lp.path(view))
		if !nestv.IsValid() || nestv.IsZero() {
			if lp.atLeaf(view) {
				details.callback(nil)
			}
		} else {
			lp.recurse(nestv.Interface(), view+1, details)
		}

	case reflect.Ptr:
		lp.recurse(reflect.ValueOf(data).Elem().Interface(), view, details)

	default:
		return nil, fmt.Errorf("unhandled case: %T", data)
	}

	return nil, nil

}

func (lp *Lenspath) traverseSlice(value any, view int, details *TraversalDetails) {
	// return []any if the array is not homogeneous (some lens gets return nil for example
	// or the map entries have different types for same keys)
	// else if array is homogeneous, return []<type> (e.g. []string)

	arr := reflect.ValueOf(value)
	if arr.Len() == 0 {
		lp.recurse(nil, view+1, details)
		return
	}

	for j := 0; j < arr.Len(); j++ {
		lp.recurse(arr.Index(j).Interface(), view+1, details)
	}
}

func (lp *Lenspath) traverseMap(value any, view int, details *TraversalDetails) (any, error) {
	key := reflect.ValueOf((lp.lens[view]))
	keyv := reflect.ValueOf(value).MapIndex(key)
	var val any
	var err error

	if !keyv.IsValid() || keyv.IsZero() {
		if !lp.atLeaf(view) {
			return nil, nil
		}

		val = details.callback(nil)
	} else {
		val, err = lp.recurse(keyv.Interface(), view+1, details)
	}

	if err != nil {
		return nil, err
	}

	if details.settable && lp.atLeaf(view) {
		reflect.ValueOf(value).SetMapIndex(key, reflect.ValueOf(val))
	}

	return nil, nil
}

func (lp *Lenspath) len() int {
	return len(lp.lens)
}

func (lp *Lenspath) path(view int) string {
	return string(lp.lens[view])
}

func (lp *Lenspath) atLeaf(view int) bool {
	return view >= lp.len()-1
}

func (lp *Lenspath) isArrBased() bool {
	for _, lensv := range lp.lens {
		if lensv == "*" {
			return true
		}
	}
	return false
}
