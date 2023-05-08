package lenspath

import (
	"fmt"
	"reflect"
)

type Lens = string

type LeafCallback = func(any, error) any // called when a leaf (of the lenspath) is reached; return value is set to the leaf (in case of setter)

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

func (lp *Lenspath) Getter(data any, callback LeafCallback) {
	lp.recurse(data, 0, &TraversalDetails{callback: callback, settable: false})
}

func (lp *Lenspath) recurse(data any, view int, details *TraversalDetails) {
	if view == lp.len() {
		details.callback(data, nil)
		return
	} else if data == nil {
		return
	}

	kind := reflect.TypeOf(data).Kind()

	switch kind {
	case reflect.Map:
		lp.traverseMap(data, view, details)

	case reflect.Slice, reflect.Array:
		if lp.path(view) == "*" {
			lp.traverseSlice(data, view, details)
		} else {
			details.callback(nil, NewInvalidLensPathErr(view, ArrayExpectedErr))
		}

	case reflect.Struct:
		nestv := reflect.ValueOf(data).FieldByName(lp.path(view))
		if !nestv.IsValid() || nestv.IsZero() {
			if lp.atLeaf(view) {
				details.callback(nil, nil)
			} else {
				return
			}
		} else {
			lp.recurse(nestv.Interface(), view+1, details)
		}

	case reflect.Ptr:
		lp.recurse(reflect.ValueOf(data).Elem().Interface(), view, details)

	default:
		details.callback(nil, fmt.Errorf("unhandled case: %T", data))
	}

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

func (lp *Lenspath) traverseMap(value any, view int, details *TraversalDetails) {
	key := reflect.ValueOf((lp.lens[view]))
	keyv := reflect.ValueOf(value).MapIndex(key)

	if !keyv.IsValid() || keyv.IsZero() {
		if view < lp.len()-1 {
			return
		}

		details.callback(nil, nil)
	} else {
		lp.recurse(keyv.Interface(), view+1, details)
	}
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
