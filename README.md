# lenspath

`lenspath` is a Golang package for traversing and updating nested data structures, such as maps, slices, structs and arrays. It uses the functional concept of lenses to compose functions that manipulate specific parts of these structures while leaving the rest unchanged. With Lenspath, you can create lenses that focus on nested fields, specific indices of a slice or array, or a range of indices, and combine them to create more complex lenses. The package provides an API for using and combining lenses, making it a useful tool for working with complex, nested data.


```golang


data := map[string]any{
		"name":   "chacha",
		"region": "India",
		"additional": map[string]any{
			"birthmark": "cut on the left hand",
			"addi": map[string]string{
				"code":     "334532",
				"landmark": "near the forest entry",
			},
        }
}

codeLp = Lenspath(["additional", "addi", "code"])
fmt.Println(housePath.Get(data));    // "334532"

codeLp.Set(data, "5A-1001");
fmt.Println(housePath.Get(data));    // "5A-1001"
```


NOTE: it works well with maps and arrays. Setting on structs is problematic as reflection required fields to be "settable". That will come soon.

# Further
- [x] add logic to get lenspath for structs
- [x] lenspath for maps (and tests)
- [x] lenspath for arrays
- [x] option to assume nil for exhausted lens path
- [x] add set capability to lenspath
- [ ] option to be able to get/set values for unexported struct fields 
- [x] ways to combine lenses
- [ ] document and give examples for various features
- [ ] test error scenarios (proper error types must be returned)
- [ ] add ci for running unit tests



## Problem

if there are several "*" in the lenspath, when doing a get, the value you get would be nested array.
traversing these might be challenging. Maybe lenspath can provide a recursive call itself, which can 
unwrap the value and provide interface to get each leaf value sequentially. 
The problem would then be that "SET" also needs to be nested similarly to make this work.
