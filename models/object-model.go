package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"sync"
)

var (
	vaultMap = &JsonVaultMap{
		m: make(map[uuid.UUID]*Object),
	}
)

// JsonVaultMap is a model for thread-safe writing/reading json values
type JsonVaultMap struct {
	mx sync.RWMutex
	m  map[uuid.UUID]*Object
}

// load gets obj by uid
func (jvm *JsonVaultMap) load(uid uuid.UUID) (*Object, bool) {
	jvm.mx.RLock()
	defer jvm.mx.RUnlock()
	val, ok := jvm.m[uid]
	return val, ok
}

// save obj by uid
func (jvm *JsonVaultMap) save(uid uuid.UUID, obj *Object) {
	jvm.mx.Lock()
	defer jvm.mx.Unlock()
	jvm.m[uid] = obj
}

// delete deletes obj by uid
func (jvm *JsonVaultMap) delete(uid uuid.UUID) {
	jvm.mx.Lock()
	defer jvm.mx.Unlock()
	delete(jvm.m, uid)
}

// export returns map values for further file saving
func (jvm *JsonVaultMap) export() map[uuid.UUID]*Object {
	jvm.mx.Lock()
	defer jvm.mx.Unlock()
	return jvm.m
}

// enrich
func (jvm *JsonVaultMap) enrich(s map[uuid.UUID]*Object) {
	jvm.mx.Lock()
	defer jvm.mx.Unlock()
	jvm.m = s
}

// Object is a model for storing json objects
type Object struct {
	Id        *uuid.UUID  `json:"id,omitempty"`
	Obj       interface{} `json:"object,omitempty"`
	CreatedOn string      `json:"createdOn,omitempty"`
	UpdatedOn string      `json:"updatedOn,omitempty"`
	DeleteOn  string      `json:"deleteOn,omitempty"`
}

// ToJSON encodes incoming obj from io.Writer
func (obj *Object) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(obj)
}

// FromJSON decodes incoming obj from io.Reader
func (obj *Object) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(obj)
}

func GetObject(uid uuid.UUID) (*Object, bool) {
	return vaultMap.load(uid)
}

func SaveObject(uid uuid.UUID, obj *Object) {
	vaultMap.save(uid, obj)
}

func DeleteObject(uid uuid.UUID) {
	vaultMap.delete(uid)
}

func ExportStorage() map[uuid.UUID]*Object {
	return vaultMap.export()
}

func ImportStorage(storage map[uuid.UUID]*Object) {
	vaultMap.enrich(storage)
}
