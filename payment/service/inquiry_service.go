package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"payment/repository"
	"payment/util"
	"time"
)

func Inquiry(res http.ResponseWriter, req *http.Request) error {

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("cannot read body")
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Decode setelah ambil raw
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

	signature := req.Header.Get("X-SIGNATURE")
	xTimestamp := req.Header.Get("X-TIMESTAMP")

	if !util.ValidateSignature(bodyBytes, xTimestamp, signature) {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "98",
			"responseDesc": "Invalid Signature",
		})
		return errors.New("signature mismatch")
	}

	billNumber := input["billNumber"]
	fromAccount := input["fromAccount"]

	//QUERY BILL
	bill, err := repository.GetByBillNumber(billNumber)
	if err != nil || bill.BillNumber == "" {
		if bill.Status == "PAID" {
			util.MapperResponse(res, http.StatusOK, map[string]string{
				"responseCode": "94",
				"responseDesc": "Bill Already Paid",
			})
			return errors.New("invalid bill")
		}
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Bill Not Found",
		})
		return errors.New("invalid bill")
	}

	paymentToken := time.Now().UnixNano() / int64(time.Millisecond)
	paymentTokenStr := fmt.Sprintf("%d", paymentToken)
	log.Printf("paymentToken: %+v\n", paymentToken)
	repository.InquiryTrx(bill.BillNumber, bill.Amount, paymentTokenStr)

	//HISTORY LOG
	actionHist := "INQUIRY BILL"
	descriptionHist := fromAccount + " INQUIRY " + bill.Amount + " TO " + bill.BillNumber
	util.HistoryLog(actionHist, descriptionHist)

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
		"transactionDetail": map[string]string{
			"billNumber":   bill.BillNumber,
			"billAmount":   bill.Amount,
			"billType":     bill.BillType,
			"description":  bill.Description,
			"status":       bill.Status,
			"paymentToken": paymentTokenStr,
		},
	})

	return nil

}
