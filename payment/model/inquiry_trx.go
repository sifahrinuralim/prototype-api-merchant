package model

type InquiryTrx struct {
	BillNumber   string `json:"billNumber"`
	Amount       string `json:"amount"`
	PaymentToken string `json:"paymentToken"`
}
