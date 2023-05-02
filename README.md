# lenspath

`lenspath` is a Golang package for traversing and updating nested data structures, such as maps, slices, structs and arrays. It uses the functional concept of lenses to compose functions that manipulate specific parts of these structures while leaving the rest unchanged. With Lenspath, you can create lenses that focus on nested fields, specific indices of a slice or array, or a range of indices, and combine them to create more complex lenses. The package provides an API for using and combining lenses, making it a useful tool for working with complex, nested data.


```golang

type User struct {
    name string
    address *Address
}

type Address struct {
    house string
    street string
    area string
    city string
    country string
}

user := User {
    name: "Sudeep",
    address: &Address {
        house: "124",
        street: "street 1",
        area: "Mango",
        city: "Jamshedpur",
        country: "India",
    }
}

housePath = Lenspath(["address", "house"])
fmt.Println(housePath.Get(user));    // "124"

housePath.Set(user, "5A-1001");
fmt.Println(housePath.Get(user));    // "5A-1001"
fmt.Println(user.address.house);     // "5A-1001"

```


# Further
- [x] add logic to get lenspath for structs
- [x] lenspath for maps (and tests)
- [x] lenspath for arrays
- [x] option to assume nil for exhausted lens path
- [ ] add set capability to lenspath
- [ ] option to be able to get/set values for unexported struct fields 
- [ ] ways to combine lenses
- [ ] document and give examples for various features