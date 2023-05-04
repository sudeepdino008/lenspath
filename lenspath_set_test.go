package lenspath

import "testing"

func TestMapSet_toplevel(t *testing.T) {
	data := map[string]string{
		"name":   "chacha",
		"region": "himalayas",
	}

	checkSetWithLensPath(t, data, []string{"name"}, "chacha_new", false)
}

func TestMapSet_internal(t *testing.T) {
	data := getTestMap()
	checkSetWithLensPath(t, data, []string{"name"}, "chacha_new", false)

	checkSetWithLensPath(t, data, []string{"additional", "birthmark"}, "2.cut on the right hand", false)

	checkSetWithLensPath(t, data, []string{"additional", "addi", "code"}, "334532_new", false)

	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_h"}, []string{"too heavy", "too light"}, false)

	// tag_w is empty for some entries
	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []string{"tag_w_new", "tag_w_new2"}, false)

	// // lenspath can't be fully traversed (should lead to error)
	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_n", "tag_n_1"}, []string{"tag_n_new", "tag_n_new2"}, true)

	checkSetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_n", "tag_n_1"}, []string{"tag_n_new", "tag_n_new2"}, false)
}
