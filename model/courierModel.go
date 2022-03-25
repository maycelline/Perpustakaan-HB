package model

type Courier struct {
	ID          int    `json:"idCourier"`
	CourierName string `json:"courierName"`
	NumberPlate string `json:"numberPlate"`
}
