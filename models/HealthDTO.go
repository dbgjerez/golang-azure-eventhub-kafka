package models

// health status
const (
	HealhStatusUp   = "UP"
	HealhStatusDown = "DOWN"
)

type HealthDTO struct {
	Status string `json:"status"`
}
