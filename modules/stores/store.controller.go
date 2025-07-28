package stores

import (
	custom_errors "digishop/utilities/errors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
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
func (s storeController) updateStoreCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
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
	isError, customErr := s.service.UpdateStoreSrv(store)
	if isError {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Store updated successfully"}`)
}
func (s storeController) getStoreByUserIDCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
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
func (s storeController) createNewProductCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
		return
	}
	err := json.Unmarshal([]byte(headerUserData), &userData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid sender data"}`)
		return
	}
	var product productRequest
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid request body"}`)
		return
	}
	validator := validator.New()
	err = validator.Struct(product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		objError := custom_errors.ParseError(err)
		strError, _ := json.Marshal(objError)
		fmt.Fprintf(w, `{"errors":%s}`, strError)
		return
	}
	product.UserID = userData["id"].(string)
	isError, customErr := s.service.CreateNewProductSrv(product)
	if isError {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message":"Product created successfully"}`)
}
func (s storeController) updateProductCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
		return
	}
	err := json.Unmarshal([]byte(headerUserData), &userData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid sender data"}`)
		return
	}
	var product updateProductRequest
	err = json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid request body"}`)
		return
	}
	validator := validator.New()
	err = validator.Struct(product)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		objError := custom_errors.ParseError(err)
		strError, _ := json.Marshal(objError)
		fmt.Fprintf(w, `{"errors":%s}`, strError)
		return
	}
	product.UserID = userData["id"].(string)
	isError, customErr := s.service.UpdateProductSrv(product)
	if isError {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Product updated successfully"}`)
}
func (s storeController) getListProductCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
		return
	}
	err := json.Unmarshal([]byte(headerUserData), &userData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid sender data"}`)
		return
	}
	query := r.URL.Query()
	var param getListProductRequest
	search := query.Get("search")
	param.Search = search
	page := query.Get("paginationPage")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		pageNum = 1 // Default to first page if conversion fails
	}
	param.PaginationPage = uint(pageNum)
	rows := query.Get("paginationRow")
	pageRow, err := strconv.Atoi(rows)
	if err != nil {
		pageRow = 20
	}
	param.PaginationRows = uint(pageRow)
	sortField := query.Get("sortField")
	param.SortField = sortField
	sortOrder := query.Get("sortOrder")
	param.SortOrder = sortOrder
	param.UserID = userData["id"].(string)
	validator := validator.New()
	err = validator.Struct(param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		objError := custom_errors.ParseError(err)
		strError, _ := json.Marshal(objError)
		fmt.Fprintf(w, `{"errors":%s}`, strError)
		return
	}
	data, customErr := s.service.GetListProductSrv(param)
	if customErr != (custom_errors.CustomError{}) {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	if data == nil {
		data = []getListProductResponse{}
	}
	strData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func (s storeController) getDetailProductCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"errors":"ID is required"}`)
		return
	}
	data, customErr := s.service.GetProductDetailSrv(id)
	if customErr != (custom_errors.CustomError{}) {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	strData, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func (s storeController) deleteStoreProductCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userData map[string]interface{}
	headerUserData := r.Header.Get("X-User-Data")
	if headerUserData == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"errors":"Unauthencticated"}`)
		return
	}
	err := json.Unmarshal([]byte(headerUserData), &userData)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid sender data"}`)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"ID is required"}`)
		return
	}
	isError, customErr := s.service.DeleteStoreProductSrv(id, userData["id"].(string))
	if isError {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message":"Product deleted successfully"}`)
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
