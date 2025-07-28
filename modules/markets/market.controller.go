package markets

import (
	custom_errors "digishop/utilities/errors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type marketController struct {
	service marketService
}

var controller marketController

func (m marketController) getAllCategoryCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	categories, customErr := m.service.GetAllCategorySrv()
	if customErr.Code != 0 {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	strData, err := json.Marshal(categories)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func (m marketController) getListProductByCategoryCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	listProduct, customErr := m.service.GetListProductByCategorySrv(r.URL.Query().Get("categoryID"))
	if customErr.Code != 0 {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	strData, err := json.Marshal(listProduct)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func (m marketController) getProductDetailByIDCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	product, customErr := m.service.GetProductDetailByIDSrv(id)
	if customErr.Code != 0 {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	strData, err := json.Marshal(product)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func (m marketController) exploreProductsCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	listProduct, customErr := m.service.ExploreProductsSrv(r.URL.Query().Get("search"))
	if customErr.Code != 0 {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	strData, err := json.Marshal(listProduct)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func (m marketController) manageCartCtrl(w http.ResponseWriter, r *http.Request) {
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
	var param manageCartRequest
	err = json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"errors":"Invalid request body"}`)
		return
	}
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
	customErr := m.service.ManageCartSrv(param.UserID, param.ProductID, param.Quantity)
	if customErr.Code != 0 {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, "Success Add to Cart")
}
func (m marketController) getUserCartCtrl(w http.ResponseWriter, r *http.Request) {
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

	data, customErr := m.service.GetUserCartsSrv(userData["id"].(string))
	if customErr.Code > 0 {
		log.Println(customErr)
		w.WriteHeader(int(customErr.Code))
		fmt.Fprintf(w, `{"errors":"%s"}`, customErr.MessageToSend)
		return
	}
	if data == nil {
		data = []cartData{}
	}
	strData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling data: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"errors":"%s"}`, "Internal server error")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"data":%s}`, strData)
}
func factoryMarketController(repo iRepo) marketController {
	if controller == (marketController{}) {
		service := factoryMarketService(repo)
		controller = marketController{service: service}
	}
	return controller
}
