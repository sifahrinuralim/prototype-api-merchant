package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"payment/model"
	"payment/repository"
	"payment/util"
)

func Payment(res http.ResponseWriter, req *http.Request) error {

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("cannot read body")
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var input map[string]string
	errRequest := json.Unmarshal(bodyBytes, &input)
	if errRequest != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	auth := req.Header.Get("Authorization")
	claims, err := util.ValidateAndExtractClaims(auth)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Unauthorized Access",
		})
		return errors.New("unauthorized token")
	}
	log.Printf("claims: %+v\n", claims)

	credential := claims["credential"].(string)
	transactionId := input["transactionId"]
	paymentToken := input["paymentToken"]
	amount := input["amount"]
	billNumber := input["billNumber"]
	fromAccount := input["fromAccount"]
	pin := input["pin"]

	//SIGNATURE VALIDATION
	signature := req.Header.Get("X-SIGNATURE")
	xTimestamp := req.Header.Get("X-TIMESTAMP")
	if !util.ValidateSignature(bodyBytes, xTimestamp, signature) {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "98",
			"responseDesc": "Invalid Signature",
		})
		return errors.New("signature mismatch")
	}

	//PIN VALIDATION
	payloadPinValidation := model.PinValidation{Credential: credential, Pin: pin}
	var pinValidationRes model.PinValidationResponse
	errPinValidation := ConnectInternal("http://localhost:8090/v1/user/pin-validation", payloadPinValidation, &pinValidationRes, auth)
	if errPinValidation != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "99",
			"responseDesc": "General Error",
		})
		return errors.New("invalid pin")
	}

	if pinValidationRes.ResponseCode != "00" {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED PIN " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": pinValidationRes.ResponseCode,
			"responseDesc": pinValidationRes.ResponseDesc,
		})
		return errors.New("invalid pin")
	}

	//CHECK BILL
	inquiryTrx, err := repository.GetByPaymentTokenAndBillNumber(paymentToken, billNumber)
	if err != nil || inquiryTrx.PaymentToken == "" {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED BILL NOT FOUND " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Bill Not Found",
		})
		return errors.New("invalid bill")
	}

	if inquiryTrx.Amount != amount {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED INVALID AMOUNT " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "13",
			"responseDesc": "Invalid Amount",
		})
		return errors.New("invalid amount")
	}
	rrn := "P001" + transactionId

	//CHECK AMOUNT CBS
	payloadAccountNo := model.AccountInfo{AccountNo: fromAccount}
	var accountNoRes model.AccountInfoResponse
	errAccountNo := ConnectInternal("http://localhost:8090/v1/account/info", payloadAccountNo, &accountNoRes, auth)
	if errAccountNo != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "99",
			"responseDesc": "General Error",
		})
		return errors.New("general error")
	}

	if accountNoRes.ResponseCode != "00" {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED AMOUNT " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": accountNoRes.ResponseCode,
			"responseDesc": accountNoRes.ResponseDesc,
		})
		return errors.New("invalid account")
	}

	if accountNoRes.TransactionDetail.IsDormant == "Y" {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED ACCOUNT DORMANT " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Account No Is Dormant",
		})
		return errors.New("invalid account")
	}

	//DEBET AMOUNT
	payloadMainTransaction := model.MainTransaction{AccountNo: fromAccount, Amount: amount}
	var mainTransactionRes model.MainTransactionResponse
	errMainTransaction := ConnectInternal("http://localhost:8090/v1/account/main-transaction", payloadMainTransaction, &mainTransactionRes, auth)
	if errMainTransaction != nil {
		fmt.Println("Error:", err)
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Invalid Account No",
		})
		return errors.New("invalid account")
	}

	if mainTransactionRes.ResponseCode != "00" {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": mainTransactionRes.ResponseCode,
			"responseDesc": mainTransactionRes.ResponseDesc,
		})
		return errors.New(mainTransactionRes.ResponseDesc)
	}

	//SAVE PAYMENT TRX
	repository.PaymentTrx(fromAccount, billNumber, amount, paymentToken, transactionId, rrn)

	//QUERY BILL
	bill, err := repository.GetByBillNumber(billNumber)
	if err == nil && bill.BillNumber == "" {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED BILL NOT FOUND " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Data Not Found",
		})
		return errors.New("invalid bill")
	}

	if bill.Status == "PAID" {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "FAILED BILL ALREADY PAID " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "94",
			"responseDesc": "Bill Already Paid",
		})
		return errors.New("invalid bill")
	} else {
		//HISTORY LOG
		actionHist := "PAYMENT BILL"
		descriptionHist := "SUCCESS " + fromAccount + " PAYMENT " + amount + " TO " + billNumber
		util.HistoryLog(actionHist, descriptionHist)

		err = repository.UpdateStatus(billNumber, "PAID")
		if err != nil {
			util.MapperResponse(res, http.StatusOK, map[string]string{
				"responseCode": "99",
				"responseDesc": "System Error",
			})
			return err
		}
	}

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
		"transactionDetail": map[string]string{
			"transactionId": transactionId,
			"rrn":           rrn,
		},
	})

	return nil

}
