package tests

import (
	"testing"

	"github.com/sudeepdino008/lenspath"
)

type TestStruct struct {
	name       string
	count      int
	additional TestStructNested
}

type TestStructNested struct {
	code string
	addi *TestStructNested
}

func TestGolangBindings(t *testing.T) {
	if lp, err := lenspath.Create([]string{"name"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v, err := lp.Get(TestStruct{name: "test"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v != "test" {
		t.Errorf("Expected \"test\", got %v", v)
	}

	if lp, err := lenspath.Create([]string{"count", "additional", "code"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v, err := lp.Get(TestStruct{count: 1, additional: TestStructNested{code: "test2"}}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	} else if v != "test2" {
		t.Errorf("Expected \"test2\", got %v", v)
	}
}
