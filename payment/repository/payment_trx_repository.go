package repository

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"payment/model"
)

func PaymentTrx(fromAccount, billNumber, amount, paymentToken, transactionId, rrn string) {
	file := "storage/payment_trx.json"
	data, _ := ioutil.ReadFile(file)
	var paymentTrx []model.PaymentTrx
	json.Unmarshal(data, &paymentTrx)

	paymentTrx = append(paymentTrx, model.PaymentTrx{
		FromAccount:   fromAccount,
		BillNumber:    billNumber,
		Amount:        amount,
		PaymentToken:  paymentToken,
		TransactionId: transactionId,
		RRN:           rrn,
	})
	newData, _ := json.MarshalIndent(paymentTrx, "", "  ")
	ioutil.WriteFile(file, newData, os.ModePerm)
}
