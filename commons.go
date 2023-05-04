package lenspath

import (
	"reflect"
	"testing"
)

func checkSetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any, setFail bool) {
	lp, err := Create(lens)
	if err != nil {
		t.Errorf("create: Expected no error, got %v", err)
		return
	}

	_, err = lp.Set(structure, expectedValue)
	if err != nil && !setFail {
		t.Errorf("set: Expected no error, got %v", err)
		return
	} else if err == nil && setFail {
		t.Errorf("set: Expected error, got %v", err)
		return
	} else if err != nil && setFail {
		// success
		return
	}

	new_value, err := lp.Get(structure)
	if err != nil {
		t.Errorf("og_get: Expected no error, got %v", err)
		return
	}

	if !reflect.DeepEqual(new_value, expectedValue) {
		t.Errorf("compare: Expected %v, got %v", expectedValue, new_value)
		return
	}
}

func checkGetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any, createFail bool, getFail bool, assumeNil bool) {
	lp, err := Create(lens)

	switch {
	case err != nil && !createFail:
		t.Errorf("Expected no error, got %v", err)

	case err == nil && createFail:
		t.Errorf("Expected error, got %v", lp)

	case err != nil && createFail:
		// success
		return
	}

	lp.WithOptions(WithAssumeNil(assumeNil))

	value, err := lp.Get(structure)

	switch {
	case err != nil && !getFail:
		t.Errorf("Expected no error, got %v", err)

	case err != nil && getFail:
		// success
		return

	case err == nil && getFail:
		t.Errorf("Expected error, got %v", value)

	case reflect.ValueOf(value).Kind() == reflect.Slice || reflect.ValueOf(value).Kind() == reflect.Array:
		if !reflect.DeepEqual(value, expectedValue) {
			t.Errorf("Expected %v, got %v", expectedValue, value)
		}

	case value != expectedValue:
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}
}
