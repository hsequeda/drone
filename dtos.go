package drone

// RegisterDroneDTO struct is the value passed in the body of POST /registerDrone.
type RegisterDroneDTO struct {
	Serial      string     `json:"serial"`
	Model       DroneModel `json:"model"`
	WeightLimit uint32     `json:"weight_limit"`
	Battery     uint8      `json:"battery"`
}

// RegisterDroneDTO struct is the value passed in the body of PUT /drone/{serial}.
type LoadMedicationDTO struct {
	Name   string `json:"name"`
	Weight uint32 `json:"weight"`
	Code   string `json:"code"`
}

// RegisterDroneDTO struct is used in the response of GET /drone/{serial}/medications
type MedicationDTO struct {
	Name   string `json:"name"`
	Weight uint32 `json:"weight"`
	Code   string `json:"code"`
	Image  string `json:"picture_path"`
}

// AvailableDroneDTO struct is used in the response of GET /drones
type AvailableDroneDTO struct {
	Serial          string     `json:"serial"`
	Model           DroneModel `json:"model"`
	WeightLimit     uint32     `json:"weight_limit"`
	BatteryCapacity uint8      `json:"battery_capacity"`
	ConsumedWeight  uint32     `json:"consumed_weight"`
	State           DroneState `json:"state"`
}

type DroneBatteryLevelDTO struct {
	BatteryLevel uint8 `json:"BatteryLevel"`
}
