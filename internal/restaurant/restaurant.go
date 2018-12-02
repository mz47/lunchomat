package restaurant

// Restaurant struct
type Restaurant struct {
	Name     string
	Distance float64
	Rating   float64
	Link     string
}

// NewRestaurant ...
func NewRestaurant(name string, distance float64, rating float64, url string) Restaurant {
	return Restaurant{name, distance, rating, url}
}
