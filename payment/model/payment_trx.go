package model

type PaymentTrx struct {
	FromAccount   string `json:"fromAccount"`
	BillNumber    string `json:"billNumber"`
	Amount        string `json:"amount"`
	PaymentToken  string `json:"paymentToken"`
	TransactionId string `json:"transactionId"`
	RRN           string `json:"rrn"`
}
