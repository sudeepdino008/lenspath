# lenspath

`lenspath` used to create a deep lens, which is a specialized function that can both retrieve and update properties in deeply nested objects. The lensPath function takes an array of property names and returns a lens that focuses on the value at that path.

You can then set or view the lenspath value.

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
- [x] add logic to get/set lenspath for structs
- [x] lenspath for maps (and tests)
- [ ] option to assume nil for exhausted lens path
- [ ] option to be able to get/set values for unexported struct fields 
