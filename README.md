# lenspath

`lenspath` is a Golang package for traversing and updating nested data structures, such as maps, slices, ~~structs~~ and arrays. It uses the functional concept of lenses to compose functions that manipulate specific parts of these structures while leaving the rest unchanged. With Lenspath, you can create lenses that focus on nested fields (including arrays) and combine them to create more complex lenses. The package provides an API for using and combining lenses, making it a useful tool for working with complex, nested data.


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
		},
	}

codePath, _ := Create([]string{"additional", "addi", "code"})
value, _ := codePath.Get(data)
fmt.Println(value) // "334532"

codePath.Set(data, "5A-1001")
value, _ = codePath.Get(data)
fmt.Println(value) // "5A-1001"
```


NOTE: it works well with maps and arrays. Setting on structs is problematic as reflection required fields to be "settable".  


## the callback API
The golang lenspath also provides a callback API for get and set:

The following will assign index value to the leaf nodes (`data`):
```go
dataLens, _ := Create([]string{"replicaEvent", "*", "data"})

index = 0
err := dataLens.Setter(map, func(value any) any {
    index=index+1
    return index    ## nth leaf on the same level
})
```

This API allows for more complex editing/retrieval of leaf node values (the simpler Get/Set are infact expressed in terms of Setter/Getter)
