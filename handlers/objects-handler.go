package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"json-objects-vault/models"
	"json-objects-vault/scheduler"
	"log"
	"net/http"
	"time"
)

// ObjectHandler is a simple handler
type ObjectHandler struct {
	l *log.Logger
}

// GetObject returns object from storage
func (o *ObjectHandler) GetObject(rw http.ResponseWriter, r *http.Request) {
	o.l.Println("Objects Handler GetObject")

	//req to id path var
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["uid"])
	if err != nil {
		http.Error(rw, "Unable to Marshal Object Id", http.StatusBadRequest)
		return
	}
	o.l.Printf("Retrieving object with id = %s", id)

	//req to data vault with id
	obj, val := models.GetObject(id)

	if val != true {
		obj = &models.Object{}
	}

	//parse data value to rw
	rw.Header().Set("Content-Type", "application/json")
	err = obj.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}

// StoreObject puts object to storage
func (o *ObjectHandler) StoreObject(rw http.ResponseWriter, r *http.Request) {
	o.l.Println("Objects Handler StoreObject")

	//req to id path var
	vars := mux.Vars(r)

	id, err := uuid.Parse(vars["uid"])

	if err != nil {
		http.Error(rw, "Unable to Marshal Object Id", http.StatusBadRequest)
		return
	}
	o.l.Printf("Updating object with id = %s", id)

	var body interface{}
	e := json.NewDecoder(r.Body)
	err = e.Decode(&body)
	if err != nil {
		http.Error(rw, "Unable to Marshal JSON", http.StatusBadRequest)
		return
	}
	o.l.Printf("Marshaled body is #%v", body)

	//get header value
	expiresTime := r.Header.Get("Expires")

	//building models.Object inst
	obj := models.Object{
		Id:        &id,
		Obj:       body,
		CreatedOn: time.Now().String(),
		DeleteOn:  expiresTime,
	}

	models.SaveObject(id, &obj)

	//add schedule task for Expires Header
	scheduler.AddTask(o.l, &obj)

}

// NewObject creates a new Handler with given logger
func NewObject(l *log.Logger) *ObjectHandler {
	return &ObjectHandler{l: l}
}
