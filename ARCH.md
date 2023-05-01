

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