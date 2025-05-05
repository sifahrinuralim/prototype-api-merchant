package repository

import (
	"encoding/json"
	"io/ioutil"
	"ws/model"
)

func GetByTransactionType(transactionType string) (model.Bridge, error) {

	data, _ := ioutil.ReadFile("storage/bridge.json")
	var listBridge []model.Bridge
	json.Unmarshal(data, &listBridge)
	for _, lb := range listBridge {
		if lb.TransactionType == transactionType {
			return lb, nil
		}
	}

	return model.Bridge{}, nil

}
