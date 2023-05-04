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

	clp2 := lp.Compose([]Lens{"f", "g"})
	compose_expected = []string{"a", "b", "c", "f", "g"}
	if !reflect.DeepEqual(clp2.lens, compose_expected) {
		t.Errorf("composed: Expected %v, got %v", compose_expected, clp2.lens)
	}

	// previous compose should hold
	compose_expected = []string{"a", "b", "c", "d", "e"}
	if !reflect.DeepEqual(clp.lens, compose_expected) {
		t.Errorf("composed: Expected %v, got %v", compose_expected, clp.lens)
	}

	clp = clp2.Compose([]Lens{"h", "i"})
	compose_expected = []string{"a", "b", "c", "f", "g", "h", "i"}
	if !reflect.DeepEqual(clp.lens, compose_expected) {
		t.Errorf("composed: Expected %v, got %v", compose_expected, clp.lens)
	}
}

func TestCompose2(t *testing.T) {
	lp, _ := Create([]string{"a", "b", "c"})
	clp := lp.Compose([]Lens{"d"})

	clp2 := clp.Compose([]Lens{"f"})
	clp3 := clp.Compose([]Lens{"g"})

	expected := []string{"a", "b", "c", "d", "f"}
	if !reflect.DeepEqual(clp2.lens, expected) {
		t.Errorf("original: Expected %v, got %v", expected, clp2.lens)
	}

	expected = []string{"a", "b", "c", "d", "g"}
	if !reflect.DeepEqual(clp3.lens, expected) {
		t.Errorf("composed: Expected %v, got %v", expected, clp3.lens)
	}
}
