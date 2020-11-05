package orders

type orderStatus string

const (
	OrderStatusUnassigned orderStatus = "UNASSIGNED"
	OrderStatusTaken      orderStatus = "TAKEN"
	OrderStatusSuccess    orderStatus = "SUCCESS"
)

func (e orderStatus) String() string {
	switch e {
	case OrderStatusUnassigned:
		return "UNASSIGNED"
	case OrderStatusTaken:
		return "TAKEN"
	case OrderStatusSuccess:
		return "SUCCESS"
	default:
		return ""
	}
}

type responseStatus string

const (
	Denied responseStatus = "REQUEST_DENIED"
)

type Order struct {
	ID       string      `json:"id"    bson:"id"`
	Distance int         `json:"distance"    bson:"distance"`
	Status   orderStatus `json:"status"    bson:"status"`
}

type Response struct {
	Status     responseStatus `json:"status"`
	StatusCode int            `json:"statusCode"`
	Rows       []Row          `json:"rows"`
}

type Row struct {
	Elements []Element `json:"elements"`
}

type Element struct {
	Distance struct {
		Value *int `json:"value"`
	} `json:"distance"`
}

type TakeOrderResponse struct {
	Status orderStatus `json:"status"`
}
