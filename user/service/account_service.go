package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"user/repository"
	"user/util"
)

func MainTransaction(res http.ResponseWriter, req *http.Request) error {
	var input map[string]string
	errPayload := json.NewDecoder(req.Body).Decode(&input)
	if errPayload != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	auth := req.Header.Get("Authorization")
	log.Printf("auth: %+v\n", auth)
	claims, err := util.ValidateAndExtractClaims(auth)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Unauthorized Access",
		})
		return errors.New("unauthorized token")
	}
	log.Printf("claims: %+v\n", claims)
	cif := claims["cif"].(string)

	accountNo := input["accountNo"]
	amountStr := input["amount"]

	account, err := repository.GetByAccountNoAndCif(accountNo, cif)
	if err != nil || account.AccountNo == "" {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Account not found",
		})
		return errors.New("account not found")
	}

	available, err1 := strconv.ParseFloat(account.AvailableAmount, 64)
	amount, err2 := strconv.ParseFloat(amountStr, 64)
	if err1 != nil || err2 != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Invalid amount format",
		})
		return errors.New("invalid amount format")
	}

	if amount > available {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "51",
			"responseDesc": "Insufficient balance",
		})
		return errors.New("insufficient balance")
	}

	// Update balance
	account.AvailableAmount = strconv.FormatFloat(available-amount, 'f', 2, 64)

	// Update the file
	err = repository.UpdateAccount(account)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "96",
			"responseDesc": "System error",
		})
		return errors.New("failed to update account")
	}

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Transaction successful",
		"newBalance":   account.AvailableAmount,
	})
	return nil
}

func AccountInfo(res http.ResponseWriter, req *http.Request) error {
	var input map[string]string
	errPayload := json.NewDecoder(req.Body).Decode(&input)
	if errPayload != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	auth := req.Header.Get("Authorization")
	log.Printf("auth: %+v\n", auth)
	claims, err := util.ValidateAndExtractClaims(auth)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Unauthorized Access",
		})
		return errors.New("unauthorized token")
	}
	log.Printf("claims: %+v\n", claims)
	cif := claims["cif"].(string)

	accountNo := input["accountNo"]

	//QUERY ACCOUNT INFO
	account, err := repository.GetByAccountNoAndCif(accountNo, cif)
	if err != nil || account.AccountNo == "" {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Account Not Found",
		})
		return errors.New("invalid account")
	}

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
		"transactionDetail": map[string]string{
			"availableAmount": account.AvailableAmount,
			"isDormant":       account.IsDormant,
		},
	})

	return nil
}
