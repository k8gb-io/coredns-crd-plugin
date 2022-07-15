package gateway

type loadbalance string

const (
	lbNone   loadbalance = "none"
	lbWeight loadbalance = "weight"
)
