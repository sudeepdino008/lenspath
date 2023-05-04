package lenspath

import (
	"reflect"
	"testing"
)

func TestCompose(t *testing.T) {
	lp, _ := Create([]string{"a", "b", "c"})
	clp := lp.Compose([]Lens{"d", "e"})

	expected := []string{"a", "b", "c"}
	//
	if !reflect.DeepEqual(lp.lens, expected) {
		t.Errorf("original: Expected %v, got %v", expected, lp.lens)
	}

	compose_expected := []string{"a", "b", "c", "d", "e"}
	if !reflect.DeepEqual(clp.lens, compose_expected) {
		t.Errorf("composed: Expected %v, got %v", compose_expected, clp.lens)
	}

	clp = lp.Compose([]Lens{"f", "g"})
	compose_expected = []string{"a", "b", "c", "f", "g"}
	if !reflect.DeepEqual(clp.lens, compose_expected) {
		t.Errorf("composed: Expected %v, got %v", compose_expected, clp.lens)
	}

	clp = clp.Compose([]Lens{"h", "i"})
	compose_expected = []string{"a", "b", "c", "f", "g", "h", "i"}
	if !reflect.DeepEqual(clp.lens, compose_expected) {
		t.Errorf("composed: Expected %v, got %v", compose_expected, clp.lens)
	}
}
