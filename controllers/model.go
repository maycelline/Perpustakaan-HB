package controllers

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

type Courier struct {
	ID          int    `json:"id"`
	CourierName string `json:"name"`
	Age         int    `json:"age"`
	Email       string `json:"email"`
}

type DataBorrowed struct {
	UserName    string `json:"username"`
	OrderDate   string `json:"order_date"`
	CourierName string `json:"courier_name"`
	Time        string `json:"time"`
}
