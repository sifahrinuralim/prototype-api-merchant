package repository

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"user/model"
)

func GetByEmailAndPassword(email, password string) (model.User, error) {

	data, _ := ioutil.ReadFile("storage/user.json")
	var listUser []model.User
	json.Unmarshal(data, &listUser)
	for _, u := range listUser {
		if (u.Email == email) && (u.Password == password) {
			return u, nil
		}
	}

	return model.User{}, nil
}

func UpdateFlag(credential, isLogin string) error {
	data, err := ioutil.ReadFile("storage/user.json")
	if err != nil {
		return err
	}

	var listUser []model.User
	json.Unmarshal(data, &listUser)

	updated := false
	for i, u := range listUser {
		if u.Credential == credential {
			listUser[i].IsLogin = isLogin
			updated = true
			break
		}
	}

	if !updated {
		return errors.New("user not found")
	}

	updatedData, err := json.MarshalIndent(listUser, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("storage/user.json", updatedData, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}
