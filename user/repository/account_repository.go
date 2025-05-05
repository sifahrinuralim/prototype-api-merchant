package repository

import (
	"encoding/json"
	"io/ioutil"
	"user/model"
)

func GetByAccountNoAndCif(accountNo, cif string) (model.Account, error) {

	data, _ := ioutil.ReadFile("storage/account.json")
	var listAccount []model.Account
	json.Unmarshal(data, &listAccount)
	for _, u := range listAccount {
		if (u.AccountNo == accountNo) && (u.Cif == cif) {
			return u, nil
		}
	}

	return model.Account{}, nil
}

func UpdateAccount(updated model.Account) error {
	data, _ := ioutil.ReadFile("storage/account.json")
	var listAccount []model.Account
	json.Unmarshal(data, &listAccount)

	for i, acc := range listAccount {
		if acc.AccountNo == updated.AccountNo {
			listAccount[i] = updated
			break
		}
	}

	updatedData, err := json.MarshalIndent(listAccount, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("storage/account.json", updatedData, 0644)
}
