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

func TestStructLensPath(t *testing.T) {
	// 1. top level struct field getting
	checkGetWithLensPath(t, TestStruct{Name: "test"}, []string{"Name"}, "test", false, false)

	// 2. nested struct field getting
	checkGetWithLensPath(t,
		TestStruct{Count: 1, Additional: TestStructNested{Code: "test2"}},
		[]string{"Additional", "Code"}, "test2", false, false)

	// 3. nested struct field getting (with pointers)
	t0 := TestStructNested{Code: "test1"}
	t1 := TestStructNested{Code: "test2", Addi: &t0}
	t2 := TestStructNested{Code: "test3", Addi: &t1}
	ts := TestStruct{Count: 1, Additional: t2}

	checkGetWithLensPath(t,
		ts,
		[]string{"Additional", "Addi", "Addi", "Code"},
		"test1", false, false)

	// 4. error expected
	checkGetWithLensPath(t, TestStruct{Name: "test"}, []string{}, "", true, false)

	// 5. lenspath not exhausted
	checkGetWithLensPath(t, ts, []string{"Additional", "Addi", "Addi", "Addi", "Code"}, "", false, true)
}

func TestMapLensPath(t *testing.T) {
	data := map[string]interface{}{
		"name":   "chacha",
		"region": "India",
		"additional": map[string]interface{}{
			"birthmark": "cut on the left hand",
			"addi": map[string]interface{}{
				"code":     "334532",
				"landmark": "near the forest entry",
			},
		},
	}

	// 1. top level map field getting
	checkGetWithLensPath(t, data, []string{"name"}, "chacha", false, false)

	// 2. nested map field getting
	checkGetWithLensPath(t, data, []string{"additional", "birthmark"}, "cut on the left hand", false, false)

	// 3. errors expected
	checkGetWithLensPath(t, data, []string{"additional", "addi", "code", "extra"}, "", false, true)
}

func checkGetWithLensPath(t *testing.T, structure interface{}, lens []string, expectedValue interface{}, createFail bool, getFail bool) {
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

	value, err := lp.Get(structure)

	switch {
	case err != nil && !getFail:
		t.Errorf("Expected no error, got %v", err)

	case err != nil && getFail:
		// success
		return

	case err == nil && getFail:
		t.Errorf("Expected error, got %v", value)

	case value != expectedValue:
		t.Errorf("Expected %v, got %v", expectedValue, value)
	}
}
