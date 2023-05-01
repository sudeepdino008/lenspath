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

func TestGolangBindings(t *testing.T) {
	if lp, err := Create([]string{"Name"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v, err := lp.Get(TestStruct{Name: "test"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v != "test" {
		t.Errorf("Expected \"test\", got %v", v)
	}

	if lp, err := Create([]string{"Additional", "Code"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v, err := lp.Get(TestStruct{Count: 1, Additional: TestStructNested{Code: "test2"}}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v != "test2" {
		t.Errorf("Expected \"test2\", got %v", v)
	}

	if lp, err := Create([]string{"Additional", "Addi", "Addi", "Code"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else {
		t0 := TestStructNested{Code: "test1"}
		t1 := TestStructNested{Code: "test2", Addi: &t0}
		t2 := TestStructNested{Code: "test3", Addi: &t1}

		if v, err := lp.Get(TestStruct{Count: 1, Additional: t2}); err != nil {
			t.Errorf("Expected no error, got %v", err)
		} else if v != "test1" {
			t.Errorf("Expected \"test1\", got %v", v)
		}
	}
}
