package drone

// Model defines the different Models of drone.
type Model int8

const (
	Lightweight Model = iota + 1
	Middleweight
	Cruiserweight
	Heavyweight
)
