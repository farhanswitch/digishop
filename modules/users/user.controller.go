package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type userController struct {
	service userService
}

var controller userController

func (u userController) RegisterUserCtrl(w http.ResponseWriter, r *http.Request) {
	var param RegisterUserRequest
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val := validator.New()
	err = val.Struct(param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	isErr, errObj := u.service.RegisterUser(param)
	if isErr {
		log.Println(errObj)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, errObj.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"User created successfully"}`)
}

func (u userController) LoginUserCtrl(w http.ResponseWriter, r *http.Request) {
	var param LoginUserRequest
	w.Header().Add("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	val := validator.New()
	err = val.Struct(param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	errObj, data := u.service.LoginUser(param)
	if errObj.Code > 10 {
		log.Println(errObj)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, errObj.MessageToSend)
		return
	}
	strData, err := json.Marshal(data)
	if err != nil {
		log.Println(errObj)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data": %s}`, strData)

}
func (u userController) CheckAuthenticationCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message":"Unauthencticated"}`)
		return
	}
	token := bearerToken[7:]
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message":"Unauthencticated"}`)
		return
	}
	newToken, errObj := u.service.CheckAuthentication(token)
	if errObj.Code > 10 {
		log.Println(errObj)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message":"%s"}`, errObj.MessageToSend)
		return
	}
	w.Header().Add("XRF-TOKEN", newToken)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Authenticated","token":"%s"}`, newToken)
}

func factoryUserController(repo iRepo) userController {
	if controller == (userController{}) {
		controller = userController{
			factoryUserService(repo),
		}
	}
	return controller
}
