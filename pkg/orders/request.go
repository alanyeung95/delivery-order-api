package orders

type PlaceOrderRequest struct {
	Origin      [2]string `json:"origin"    bson:"origin" validate:"eq=2,dive,required,latitude|longitude"`
	Destination [2]string `json:"destination"    bson:"destination"  validate:"eq=2,dive,required,latitude|longitude"`
}
