package stores

import (
	custom_errors "digishop/utilities/errors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type storeController struct {
	service storeService
}

var controller storeController

func (s storeController) registerStoreCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message":"Unauthencticated"}`)
		return
	}
	err := json.Unmarshal([]byte(headerUserData), &userData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid sender data"}`)
		return
	}
	var store storeData
	err = json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid request body"}`)
		return
	}
	validator := validator.New()
	err = validator.Struct(store)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		objError := custom_errors.ParseError(err)
		strError, _ := json.Marshal(objError)
		fmt.Fprintf(w, `{"errors":%s}`, strError)
		return
	}
	store.UserID = userData["id"].(string)
	isError, customErr := s.service.RegisterStoreSrv(store)
	if isError {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Store created successfully"}`)

}
func (s storeController) getStoreByUserIDCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"message":"Unauthencticated"}`)
		return
	}
	err := json.Unmarshal([]byte(headerUserData), &userData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid sender data"}`)
		return
	}
	store, customErr := s.service.GetStoreByUserIDSrv(userData["id"].(string))
	if customErr != (custom_errors.CustomError{}) {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	strStore, _ := json.Marshal(store)
	fmt.Fprintf(w, `{"data":%s}`, strStore)
}
func factoryStoreController(repo iRepo) storeController {
	if controller == (storeController{}) {
		service := factoryStoreService(repo)
		controller = storeController{
			service: service,
		}
	}
	return controller
}
