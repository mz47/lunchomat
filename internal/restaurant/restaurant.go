package restaurant

// Restaurant struct
type Restaurant struct {
	Name     string
	Distance string
	Rating   string
}

// NewRestaurant ...
func NewRestaurant(name string) Restaurant {
	return Restaurant{name, "distance", "rating"}
}

// Visit ...
func (r *Restaurant) Visit() {
	println("visited ", r.Name)
}

// SetDistance ...
func (r *Restaurant) SetDistance(distance string) {
	r.Distance = distance
}

// SetRating ...
func (r *Restaurant) SetRating(rating string) {
	r.Rating = rating
}
