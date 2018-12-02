package restaurant

// Restaurant struct
type Restaurant struct {
	Name string
}

// NewRestaurant ...
func NewRestaurant(name string) Restaurant {
	return Restaurant{name}
}

// Visit ...
func (r *Restaurant) Visit() {
	println("visited ", r.Name)
}
