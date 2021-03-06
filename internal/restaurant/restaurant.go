package restaurant

// Restaurant struct
type Restaurant struct {
	Id           string
	Name         string
	Distance     float64
	Rating       float64
	Link         string
	TimesVisited int64
	Preferred    bool
	Ignored      bool
}

// NewRestaurant ...
func NewRestaurant(id string, name string, distance float64, rating float64, url string) Restaurant {
	return Restaurant{
		id,
		name,
		distance,
		rating,
		url,
		0,
		false,
		false,
	}
}
