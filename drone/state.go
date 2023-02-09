package drone

// State defines the different state availables in a drone.
type State int8

const (
	Idle State = iota + 1
	Loading
	Loaded
	Delivering
	Delivered
	Returning
)
