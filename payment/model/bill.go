package model

type Bill struct {
	Description string `json:"description"`
	BillNumber  string `json:"billNumber"`
	BillType    string `json:"billType"`
	Amount      string `json:"amount"`
	Status      string `json:"status"`
}
