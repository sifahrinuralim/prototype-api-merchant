package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"ws/repository"
)

func writeJSON(res http.ResponseWriter, statusCode int, data interface{}) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	json.NewEncoder(res).Encode(data)
}

func Bridge(res http.ResponseWriter, req *http.Request) error {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("cannot read body")
	}
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Parse hanya dua field: transactionType dan transactionDetail (preserve raw)
	var input map[string]json.RawMessage
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	// Ambil transactionType
	var transactionType string
	if v, ok := input["transactionType"]; ok {
		err := json.Unmarshal(v, &transactionType)
		if err != nil || transactionType == "" {
			writeJSON(res, http.StatusOK, map[string]string{
				"responseCode": "30",
				"responseDesc": "Format Error [transactionType is mandatory]",
			})
			return errors.New("invalid transactionType")
		}
	} else {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error [transactionType is missing]",
		})
		return errors.New("transactionType missing")
	}

	log.Printf("transactionType: %+v\n", transactionType)

	// Ambil transactionDetail raw (tanpa ubah)
	transactionDetailRaw, ok := input["transactionDetail"]
	if !ok {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error [transactionDetail missing]",
		})
		return errors.New("transactionDetail missing")
	}

	dataBridge, _ := repository.GetByTransactionType(transactionType)
	if dataBridge.TransactionType == "" {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error [transactionType is invalid]",
		})
		return errors.New("transaction type not found")
	}
	log.Printf("dataBridge: %+v\n", dataBridge)
	log.Printf("transactionDetail raw: %s", string(transactionDetailRaw))

	// Kirim raw JSON transactionDetail ke internal service
	client := &http.Client{}
	apiReq, err := http.NewRequest("POST", dataBridge.URL, bytes.NewBuffer(transactionDetailRaw))
	if err != nil {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return err
	}

	auth := req.Header.Get("Authorization")
	xSignature := req.Header.Get("X-SIGNATURE")
	xTimestamp := req.Header.Get("X-TIMESTAMP")

	apiReq.Header.Set("Content-Type", "application/json")
	apiReq.Header.Set("Authorization", auth)
	apiReq.Header.Set("X-SIGNATURE", xSignature)
	apiReq.Header.Set("X-TIMESTAMP", xTimestamp)

	apiRes, err := client.Do(apiReq)
	if err != nil {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "99",
			"responseDesc": "Failed to call external API",
		})
		return err
	}
	defer apiRes.Body.Close()

	body, err := io.ReadAll(apiRes.Body)
	if err != nil {
		writeJSON(res, http.StatusOK, map[string]string{
			"responseCode": "99",
			"responseDesc": "Failed to read external response",
		})
		return err
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(apiRes.StatusCode)
	res.Write(body)

	return nil
}
