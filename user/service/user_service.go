package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"user/config"
	"user/repository"
	"user/util"

	"golang.org/x/crypto/bcrypt"
)

func HashPIN(pin string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPIN(hashedPIN, inputPIN string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPIN), []byte(inputPIN))
	return err == nil
}

func PinValidation(res http.ResponseWriter, req *http.Request) error {

	auth := req.Header.Get("Authorization")
	_, err := util.ValidateAndExtractClaims(auth)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Unauthorized Access",
		})
		return errors.New("unauthorized token")
	}

	var input map[string]string
	errPayload := json.NewDecoder(req.Body).Decode(&input)
	if errPayload != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	credential := input["credential"]
	pin := input["pin"]

	hsmPin, err := repository.GetByCredential(credential)
	if err != nil || hsmPin.Credential == "" {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Invalid PIN",
		})
		return errors.New("invalid pin")
	}

	isValid := VerifyPIN(hsmPin.Pin, pin)
	if !isValid {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Invalid PIN",
		})
		return errors.New("invalid pin")
	}

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
	})

	return nil
}

func Logout(res http.ResponseWriter, req *http.Request) error {

	auth := req.Header.Get("Authorization")
	claims, err := util.ValidateAndExtractClaims(auth)
	if err != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Unauthorized Access",
		})
		return errors.New("unauthorized token")
	}

	credential := claims["credential"].(string)

	//UPDATE IS LOGIN
	if err := repository.UpdateFlag(credential, "Y"); err != nil {
		log.Printf("failed to update login status: %v", err)
	}

	//HISTORY LOG
	actionHist := "LOGOUT"
	descriptionHist := "SUCCESS CREDENTIAL " + credential
	util.HistoryLog(actionHist, descriptionHist)

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
	})

	return nil

}

func Login(res http.ResponseWriter, req *http.Request) error {

	var input map[string]string
	errPayload := json.NewDecoder(req.Body).Decode(&input)
	if errPayload != nil {
		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "30",
			"responseDesc": "Format Error",
		})
		return errors.New("invalid payload")
	}

	email := input["email"]
	password := input["password"]
	auth := req.Header.Get("Authorization")

	log.Printf("auth: %+v\n", auth)

	if auth != config.StaticToken {
		//HISTORY LOG
		actionHist := "LOGIN"
		descriptionHist := "FAILED INVALID AUTHORIZATION" + email
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "55",
			"responseDesc": "Unauthorized Access",
		})
		return errors.New("unauthorized token")
	}

	user, err := repository.GetByEmailAndPassword(email, password)
	if err != nil || user.Email == "" {
		//HISTORY LOG
		actionHist := "LOGIN"
		descriptionHist := "FAILED INVALID EMAIL/PASSWORD" + email
		util.HistoryLog(actionHist, descriptionHist)

		util.MapperResponse(res, http.StatusOK, map[string]string{
			"responseCode": "76",
			"responseDesc": "Invalid Email or Password",
		})
		return errors.New("invalid credentials")
	}

	//GENERATE JWT TOKEN
	jwtLogin, err := util.GenerateJWT(user.Credential, user.Cif)
	log.Printf("err: %+v\n", err)

	//UPDATE IS LOGIN
	if err := repository.UpdateFlag(user.Credential, "Y"); err != nil {
		log.Printf("failed to update login status: %v", err)
	}

	//HISTORY LOG
	actionHist := "LOGIN"
	descriptionHist := "SUCCESS " + email + " CREDENTIAL " + user.Credential
	util.HistoryLog(actionHist, descriptionHist)

	util.MapperResponse(res, http.StatusOK, map[string]interface{}{
		"responseCode": "00",
		"responseDesc": "Success",
		"transactionDetail": map[string]string{
			"credential": user.Credential,
			"jwtLogin":   jwtLogin,
		},
	})

	return nil
}
