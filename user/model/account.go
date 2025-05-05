package model

type Account struct {
	Cif             string `json:"cif"`
	AccountNo       string `json:"accountNo"`
	AvailableAmount string `json:"availableAmount"`
	IsDormant       string `json:"isDormant"`
}
