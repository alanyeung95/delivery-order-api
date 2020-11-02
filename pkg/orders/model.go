package orders

type orderStatus string

const (
	Unassigned orderStatus = "UNASSIGNED"
	Taken      orderStatus = "TAKEN"
	Success    orderStatus = "SUCCESS"
)

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
		Value int `json:"value"`
	} `json:"distance"`
}
