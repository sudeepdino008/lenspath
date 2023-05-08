package lenspath

import (
	"reflect"
	"testing"
)

// func checkSetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any, setFail bool) {
// 	lp, err := Create(lens)
// 	if err != nil {
// 		t.Errorf("create: Expected no error, got %v", err)
// 		return
// 	}

// 	_, err = lp.Set(structure, expectedValue)
// 	if err != nil && !setFail {
// 		t.Errorf("set: Expected no error, got %v", err)
// 		return
// 	} else if err == nil && setFail {
// 		t.Errorf("set: Expected error, got %v", err)
// 		return
// 	} else if err != nil && setFail {
// 		// success
// 		return
// 	}

// 	new_value, err := lp.Get(structure)
// 	if err != nil {
// 		t.Errorf("og_get: Expected no error, got %v", err)
// 		return
// 	}

// 	if !reflect.DeepEqual(new_value, expectedValue) {
// 		t.Errorf("compare: Expected %v, got %v", expectedValue, new_value)
// 		return
// 	}
// }

func checkGetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any, createFail bool, getFail bool, assumeNil bool) {
	lp, err := Create(lens)

	switch {
	case err != nil && !createFail:
		t.Fatalf("Expected no error, got %v", err)

	case err == nil && createFail:
		t.Fatalf("Expected error, got %v", lp)

	case err != nil && createFail:
		// success
		return
	}

	lp.WithOptions(WithAssumeNil(assumeNil))

	containsArr := false
	for _, lensv := range lens {
		if lensv == "*" {
			containsArr = true
			break
		}
	}

	var exec_err error
	index := 0

	lp.Getter(structure, func(value any, err error) any {
		if err != nil {
			exec_err = err
		}

		if !containsArr {
			comparev(t, value, expectedValue)
		} else {
			comparev(t, value, expectedValue.([]any)[index])
			index++
		}

		return nil
	})

	if exec_err != nil && !getFail {
		t.Fatalf("Expected no error, got %v", exec_err)
	} else if exec_err == nil && getFail {
		t.Fatalf("Expected error, got %v", exec_err)
	} else if containsArr && index != len(expectedValue.([]any)) {
		t.Fatalf("expected array size mismatch")
	}

	// success

}

func comparev(t *testing.T, value any, expectedValue any) {
	if value == nil {
		if expectedValue != nil {
			t.Fatalf("Expected %v, got %v", expectedValue, value)
		} else {
			return
		}
	}
	kind := reflect.TypeOf(value).Kind()
	switch {
	case kind == reflect.Slice || kind == reflect.Array:
		if !reflect.DeepEqual(value, expectedValue) {
			t.Fatalf("Expected %v, got %v", expectedValue, value)
		}

	case value != expectedValue:
		t.Fatalf("Expected %v, got %v", expectedValue, value)
	default:
		// fine
	}
}
