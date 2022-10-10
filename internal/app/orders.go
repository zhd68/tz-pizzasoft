package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zhd68/tz-pizzasoft/internal/model"
	"github.com/zhd68/tz-pizzasoft/internal/storage"
)

func (a *App) newOrder(w http.ResponseWriter, r *http.Request) {
	order := &model.Order{}

	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := []byte(`{"error": "request body is empty"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	err = json.Unmarshal(body, order)
	if err != nil {
		msg := []byte(`{"error": "invalid 'items'"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	if err = order.Items.ValidateItems(); err != nil {
		msg := []byte(`{"error": "invalid 'items'"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	order, err = a.store.SaveOrder(order)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	data, err := json.Marshal(order)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	a.newResponse(w, http.StatusOK, data)
}

func (a *App) addItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderId := vars["order_id"]

	var items model.Items
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := []byte(`{"error": "request body is empty"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	err = json.Unmarshal(body, &items)
	if err != nil {
		msg := []byte(`{"error": "invalid 'items'"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	if err = items.ValidateItems(); err != nil {
		msg := []byte(`{"error": "invalid 'items'"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	order, err := a.store.UpdateOrder(orderId, items)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusOK, msg)
		return
	}
	data, err := json.Marshal(order)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	a.newResponse(w, http.StatusOK, data)
}

func (a *App) getOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	orderId := vars["order_id"]

	order, err := a.store.GetOrder(orderId)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusOK, msg)
		return
	}
	data, err := json.Marshal(order)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	a.newResponse(w, http.StatusOK, data)
}

func (a *App) doneOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	xAuthKey := r.Header.Get("X-Auth-Key")
	if xAuthKey != a.config.XAuthKey {
		msg := []byte(`{"error": "invalid header 'X-Auth-Key'"}`)
		a.newResponse(w, http.StatusOK, msg)
		return
	}

	vars := mux.Vars(r)
	orderId := vars["order_id"]

	order, err := a.store.DoneOrder(orderId)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusOK, msg)
		return
	}
	data, err := json.Marshal(order)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	a.newResponse(w, http.StatusOK, data)
}

func (a *App) getALLOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	xAuthKey := r.Header.Get("X-Auth-Key")
	if xAuthKey != a.config.XAuthKey {
		msg := []byte(`{"error": "invalid header 'X-Auth-Key'"}`)
		a.newResponse(w, http.StatusOK, msg)
		return
	}

	done, err := ValidateDone(r.FormValue("done"))
	if err != nil {
		msg := []byte(`{"error": "` + err.Error() + `"}`)
		a.newErrorResponse(w, err, http.StatusBadRequest, msg)
		return
	}

	orders, err := a.store.GetAllOrders(done)
	if err == storage.ErrNoSavedOrders {
		msg := []byte(`{"error": "` + err.Error() + `"}`)
		a.newResponse(w, http.StatusOK, msg)
		return
	}
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	data, err := json.Marshal(orders)
	if err != nil {
		msg := []byte(`{"error": "internal server error"}`)
		a.newErrorResponse(w, err, http.StatusInternalServerError, msg)
		return
	}
	a.newResponse(w, http.StatusOK, data)
}

func ValidateDone(done string) (*bool, error) {
	switch done {
	case "":
		return nil, nil
	case "1":
		val := new(bool)
		*val = true
		return val, nil
	case "0":
		val := new(bool)
		*val = false
		return val, nil
	default:
		return nil, fmt.Errorf("invalid done value")
	}
}
