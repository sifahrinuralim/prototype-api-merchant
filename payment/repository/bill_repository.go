package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"payment/model"
)

func GetByBillNumber(billNumber string) (model.Bill, error) {

	data, _ := ioutil.ReadFile("storage/bill.json")
	var listBill []model.Bill
	json.Unmarshal(data, &listBill)
	for _, u := range listBill {
		if u.BillNumber == billNumber {
			return u, nil
		}
	}

	return model.Bill{}, nil
}

func UpdateStatus(billNumber, status string) error {
	data, err := ioutil.ReadFile("storage/bill.json")
	if err != nil {
		return err
	}
	var listBill []model.Bill
	err = json.Unmarshal(data, &listBill)
	if err != nil {
		return err
	}
	for i, u := range listBill {
		if u.BillNumber == billNumber {
			listBill[i].Status = status
			updatedData, err := json.Marshal(listBill)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile("storage/bill.json", updatedData, 0644)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return errors.New("bill not found")
}
