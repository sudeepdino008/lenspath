package lenspath

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name       string
	Count      int
	Additional TestStructNested
}

type TestStructNested struct {
	Code string
	Addi *TestStructNested
}

func TestStructLensPath(t *testing.T) {
	// 1. top level struct field getting
	checkGetWithLensPath(t, TestStruct{Name: "test"}, []string{"Name"}, "test", false, false, false)

	// 2. nested struct field getting
	checkGetWithLensPath(t,
		TestStruct{Count: 1, Additional: TestStructNested{Code: "test2"}},
		[]string{"Additional", "Code"}, "test2", false, false, false)

	// 3. nested struct field getting (with pointers)
	t0 := TestStructNested{Code: "test1"}
	t1 := TestStructNested{Code: "test2", Addi: &t0}
	t2 := TestStructNested{Code: "test3", Addi: &t1}
	ts := TestStruct{Count: 1, Additional: t2}

	checkGetWithLensPath(t,
		ts,
		[]string{"Additional", "Addi", "Addi", "Code"},
		"test1", false, false, false)

	// 4. error expected
	checkGetWithLensPath(t, TestStruct{Name: "test"}, []string{}, "", true, false, false)

	// 5. lenspath not exhausted
	checkGetWithLensPath(t, ts, []string{"Additional", "Addi", "Addi", "Addi", "Code"}, "", false, true, false)
	// 5.1 with assumeNil
	checkGetWithLensPath(t, ts, []string{"Additional", "Addi", "Addi", "Addi", "Code"}, nil, false, false, true)
}

func TestMapLensPath(t *testing.T) {
	tagsList := []map[string]any{
		{"tag_h": "medium"},
		{"tag_w": "heavy", "tag_h": "tall"},
	}
	data := map[string]any{
		"name":   "chacha",
		"region": "India",
		"additional": map[string]any{
			"birthmark": "cut on the left hand",
			"addi": map[string]any{
				"code":     "334532",
				"landmark": "near the forest entry",
			},
			"tagsList": tagsList,
		},
	}

	// 1. top level map field getting
	checkGetWithLensPath(t, data, []string{"name"}, "chacha", false, false, false)

	// 2. nested map field getting
	checkGetWithLensPath(t, data, []string{"additional", "birthmark"}, "cut on the left hand", false, false, false)

	// 3. get all array field
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_h"}, []any{"medium", "tall"}, false, false, false)

	// 4 get all array field (but not elements have the queried nested field)
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []any{}, false, true, false)
	// 4.1 with assumeNil
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []any{nil, "heavy"}, false, false, true)

	// 4. array field getting error
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "not_found"}, nil, false, true, false)

	// 5. errors expected
	checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, "", false, true, false)
	// 5.1 with assumeNil
	checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, nil, false, true, true)

	checkGetWithLensPath(t, data, []string{"additional", "addi", "nonexisting", "extra"}, nil, false, false, true)
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
