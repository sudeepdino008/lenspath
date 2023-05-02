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

func TestStructGet(t *testing.T) {
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

func TestMapGet(t *testing.T) {
	tagsList := []map[string]any{
		{"tag_h": "medium"},
		{"tag_w": "heavy", "tag_h": "tall"},
	}
	data := map[string]any{
		"name":   "chacha",
		"region": "India",
		"additional": map[string]any{
			"birthmark": "cut on the left hand",
			"addi": map[string]string{
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
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_h"}, []string{"medium", "tall"}, false, false, false)

	// 4 get all array field (but not elements have the queried nested field)
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []string{}, false, true, false)
	// 4.1 with assumeNil
	// TODO: * lens should use closures over set methods; otherwise "absent" values will be attempted to set to nil
	// which is something that cannot be set to "string" type here.
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []string{"", "heavy"}, false, false, true)

	// 4. array field getting error
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "not_found"}, nil, false, true, false)

	// 5. errors expected
	checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, "", false, true, false)
	// 5.1 with assumeNil
	checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, nil, false, true, true)

	checkGetWithLensPath(t, data, []string{"additional", "addi", "nonexisting", "extra"}, nil, false, false, true)
}

func TestMapSet(t *testing.T) {
	data := getTestMap()
	checkSetWithLensPath(t, data, []string{"name"}, "chacha_new")

	checkSetWithLensPath(t, data, []string{"additional", "birthmark"}, "2.cut on the right hand")

	checkSetWithLensPath(t, data, []string{"additional", "addi", "code"}, "334532_new")

	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_h"}, []string{"too heavy", "too light"})
}

func TestMapSet2(t *testing.T) {
	data := getTestMap2()

	checkSetWithLensPath(t, data, []string{"name"}, "chacha_new")
}

func getTestMap() map[string]any {
	tagsList := []map[string]string{
		{"tag_h": "medium"},
		{"tag_w": "heavy", "tag_h": "tall"},
	}
	data := map[string]any{
		"name":   "chacha",
		"region": "India",
		"additional": map[string]any{
			"birthmark": "cut on the left hand",
			"addi": map[string]string{
				"code":     "334532",
				"landmark": "near the forest entry",
			},
			"tagsList":  tagsList,
			"tagsList2": tagsList,
		},
	}

	return data
}

func getTestMap2() map[string]string {
	return map[string]string{
		"name":   "chacha",
		"region": "himalayas",
	}
}

func checkSetWithLensPath(t *testing.T, structure any, lens []string, expectedValue any) {
	lp, err := Create(lens)
	if err != nil {
		t.Errorf("create: Expected no error, got %v", err)
		return
	}

	_, err = lp.Set(structure, expectedValue)
	if err != nil {
		t.Errorf("set: Expected no error, got %v", err)
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
