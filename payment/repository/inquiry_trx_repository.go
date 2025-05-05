package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"payment/model"
)

func GetByPaymentTokenAndBillNumber(paymentToken, billNumber string) (model.InquiryTrx, error) {

	data, _ := ioutil.ReadFile("storage/inquiry_trx.json")
	var listInquiryTrx []model.InquiryTrx
	json.Unmarshal(data, &listInquiryTrx)
	for _, u := range listInquiryTrx {
		if (u.PaymentToken == paymentToken) && (u.BillNumber == billNumber) {
			return u, nil
		}
	}

	return model.InquiryTrx{}, nil
}

func InquiryTrx(billNumber, amount, paymentToken string) {
	file := "storage/inquiry_trx.json"
	data, _ := ioutil.ReadFile(file)
	var inquiryTrx []model.InquiryTrx
	json.Unmarshal(data, &inquiryTrx)

	inquiryTrx = append(inquiryTrx, model.InquiryTrx{
		BillNumber:   billNumber,
		Amount:       amount,
		PaymentToken: paymentToken,
	})
	newData, _ := json.MarshalIndent(inquiryTrx, "", "  ")
	ioutil.WriteFile(file, newData, os.ModePerm)
}
