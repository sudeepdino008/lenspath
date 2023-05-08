package lenspath

// import "testing"

// func TestMapSet_toplevel(t *testing.T) {
// 	data := map[string]string{
// 		"name":   "chacha",
// 		"region": "himalayas",
// 	}

// 	checkSetWithLensPath(t, data, []string{"name"}, "chacha_new", false)
// }

// func TestMapSet_internal(t *testing.T) {
// 	data := getTestMap()
// 	checkSetWithLensPath(t, data, []string{"name"}, "chacha_new", false)

// 	checkSetWithLensPath(t, data, []string{"additional", "birthmark"}, "2.cut on the right hand", false)

// 	checkSetWithLensPath(t, data, []string{"additional", "addi", "code"}, "334532_new", false)

// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_h"}, []string{"too heavy", "too light"}, false)

// 	// tag_w is empty for some entries
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_w"}, []string{"tag_w_new", "tag_w_new2"}, false)

// 	// // lenspath can't be fully traversed (should lead to error)
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList", "*", "tag_n", "tag_n_1"}, []string{"tag_n_new", "tag_n_new2"}, true)

// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_n", "tag_n_1"}, []string{"tag_n_new", "tag_n_new2"}, false)
// }

// func TestMapSet_complex(t *testing.T) {
// 	data := getTestMap()

// 	// setting to array
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList3"}, []string{"too heavy", "too light"}, false)

// 	// []any
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList3"}, []any{nil, "too light"}, false)

// 	// setting to map
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList3"}, map[string]string{"tag_h": "too heavy", "tag_w": "too light"}, false)

// 	// // setting to struct
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList3"}, struct {
// 		Tag_h string
// 	}{Tag_h: "too heavy"}, false)

// 	// set to nil
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList3"}, nil, false)
// }

// func TestMapSet_array(t *testing.T) {
// 	data := getTestMap()

// 	// set array to array of diff size
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList4"}, []string{"a", "b", "c", "d"}, false)

// 	// // size mismatch error
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_h"}, []string{"tag_h1", "tag_h2", "tag_h3"}, true)

// 	// correct size
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList2", "*", "tag_h"}, []string{"tag_h1", "tag_h2"}, false)

// 	// setting with flattening array
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList5", "*", "tag_n", "tag_n", "*", "tag_n_1"}, []int{3, 4}, false)

// 	// different size should fail
// 	checkSetWithLensPath(t, data, []string{"additional", "tagsList5", "*", "tag_n", "tag_n", "*", "tag_n_1"}, []int{3, 4, 5}, true)
// }
