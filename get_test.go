package lenspath

import (
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

	// 5.1 lenspath not exhausted
	checkGetWithLensPath(t, ts, []string{"Additional", "Addi", "Addi", "Addi", "Code"}, nil, false, false, true)
}

func TestMapGet(t *testing.T) {
	data := getTestMap()

	// 1. top level map field getting
	checkGetWithLensPath(t, data, []string{"name"}, "chacha", false, false, false)

	// 2. nested map field getting
	checkGetWithLensPath(t, data, []string{"additional", "birthmark"}, "cut on the left hand", false, false, false)

	// 3. get all array field
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_h"}, []any{"medium", "tall"}, false, false, false)

	// 4 get all array field (but not elements have the queried nested field)
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []any{nil, "heavy"}, false, false, true)
	// 4.1 with assumeNil
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []any{nil, "heavy"}, false, false, true)

	// 4. array field getting error
	checkGetWithLensPath(t, data, []string{"additional", "tagsList", "not_found"}, nil, false, true, false)

	// 5. errors expected
	//checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, "", false, true, false)
	// 5.1 with assumeNil
	checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, nil, false, true, true)

	checkGetWithLensPath(t, data, []string{"additional", "addi", "nonexisting", "extra"}, nil, false, false, true)

	checkGetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_n", "tag_n_1"}, []any{1, nil}, false, false, true)

	checkGetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_n", "tag_n_2"}, []any{"2-string", "2-string"}, false, false, true)

	checkGetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_n", "tag_n_3"}, []any{3.0, "3.0-string"}, false, false, true)

	//	the leaves of the tree should be reported in a flattened array
	checkGetWithLensPath(t, data, []string{"additional", "tagsList5", "*", "tag_n", "tag_n", "*", "tag_n_1"}, []any{1, 3}, false, false, true)
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
			"tagsList": tagsList,
			"tagsList2": []map[string]any{
				{"tag_h": "medium", "tag_n": map[string]any{"tag_n_1": 1, "tag_n_2": "2-string", "tag_n_3": 3.0}},
				{"tag_w": "heavy", "tag_h": "tall", "tag_n": map[string]any{"tag_n_2": "2-string", "tag_n_3": "3.0-string"}},
			},
			"tagsList5": []map[string]any{
				{"tag_h": "medium", "tag_n": map[string]any{"tag_n": []map[string]any{{"tag_n_1": 1, "tag_n_2": "2-string", "tag_n_3": 3.0}}}},
				{"tag_w": "heavy", "tag_h": "tall", "tag_n": map[string]any{"tag_n": []map[string]any{{"tag_n_1": 3, "tag_n_2": "2-string", "tag_n_3": 3.0}}}},
			},
			"tagsList3": "hello world",
			"tagsList4": []string{"hello", "world", "hw"},
		},
	}

	return data
}
