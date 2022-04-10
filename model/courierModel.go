package model

type Courier struct {
	ID          int    `json:"idCourier,omitempty"`
	CourierName string `json:"courierName,omitempty"`
	NumberPlate string `json:"numberPlate,omitempty"`
}
