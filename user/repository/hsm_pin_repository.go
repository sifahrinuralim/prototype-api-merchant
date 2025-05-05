package repository

import (
	"encoding/json"
	"io/ioutil"
	"user/model"
)

func GetByCredential(credential string) (model.HsmPin, error) {

	data, _ := ioutil.ReadFile("storage/hsm_pin.json")
	var listHsm []model.HsmPin
	json.Unmarshal(data, &listHsm)
	for _, u := range listHsm {
		if u.Credential == credential {
			return u, nil
		}
	}

	return model.HsmPin{}, nil
}
