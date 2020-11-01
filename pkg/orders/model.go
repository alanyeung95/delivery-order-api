package orders

type orderStatus string

const (
	Unassigned orderStatus = "UNASSIGNED"
	Taken      orderStatus = "TAKEN"
	Success    orderStatus = "SUCCESS"
)

type Order struct {
	ID       string      `json:"id"    bson:"id"`
	Distance int         `json:"distance"    bson:"distance"`
	Status   orderStatus `json:"status"    bson:"status"`
}
