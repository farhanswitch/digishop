package markets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
func factoryMarketController(repo iRepo) marketController {
	if controller == (marketController{}) {
		service := factoryMarketService(repo)
		controller = marketController{
			service: service,
		}
	}
	return controller
}
