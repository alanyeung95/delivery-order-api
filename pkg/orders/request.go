package orders

type PlaceOrderRequest struct {
	Origin      [2]string `json:"origin"    bson:"origin"`
	Destination [2]string `json:"destination"    bson:"destination"`
}
