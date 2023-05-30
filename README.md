# lenspath

`lenspath` is a Golang package for traversing and updating nested data structures, such as maps, slices, ~~structs~~ and arrays. It uses the functional concept of lenses to compose functions that manipulate specific parts of these structures while leaving the rest unchanged. With Lenspath, you can create lenses that focus on nested fields, specific indices of a slice or array, or a range of indices, and combine them to create more complex lenses. The package provides an API for using and combining lenses, making it a useful tool for working with complex, nested data.


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


NOTE: it works well with maps and arrays. Setting on structs is problematic as reflection required fields to be "settable".  
