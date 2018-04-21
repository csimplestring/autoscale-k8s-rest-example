package customer

type Customer struct {
	Name *string  `json:"name"`
	Address *string `json:"address"`
}