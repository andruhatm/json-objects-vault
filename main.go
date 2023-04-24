package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"json-objects-vault/handlers"
	"json-objects-vault/models"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	//init logger
	l := log.New(os.Stdout, "json-vault", log.LstdFlags)

	//import existing vault sample
	initStorage(l)

	//create the handler
	objH := handlers.NewObject(l)

	//create a new serve mux and register the handlers
	sm := mux.NewRouter()

	objectsRouter := sm.Methods(http.MethodPut, http.MethodGet).Subrouter()
	objectsRouter.HandleFunc("/objects/{uid}", objH.GetObject).Methods(http.MethodGet)
	objectsRouter.HandleFunc("/objects/{uid}", objH.StoreObject).Methods(http.MethodPut)

	probesRouter := sm.Methods(http.MethodGet).Subrouter()
	probesRouter.HandleFunc("/probes/readiness", func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write([]byte("OK"))
		if err != nil {
			l.Printf("Error while writing the data to an HTTP reply with err=%s", err)
			return
		}
	})
	probesRouter.HandleFunc("/probes/liveness", func(rw http.ResponseWriter, r *http.Request) {
		//TODO check if we can access resources
		_, err := rw.Write([]byte("OK"))
		if err != nil {
			l.Printf("Error while writing the data to an HTTP reply with err=%s", err)
			return
		}
	})

	//prometheus endpoint
	sm.Handle("/metrics", promhttp.Handler())

	//http.Server instance
	s := &http.Server{
		Addr:         ":8081",
		Handler:      sm,
		TLSConfig:    nil,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port 8081")

		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	//trap os.Signal and gracefully shutdown the server
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Notify(sigCh, os.Kill)

	sig := <-sigCh
	l.Printf("Graceful shutdown with signal %s \n", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	saveStorage(l)
	s.Shutdown(ctx)

}

// initStorage imports file to a map
func initStorage(l *log.Logger) {

	f, err := os.Open("import/data.json")

	if err != nil {
		l.Printf("Error open import file err=%s \n", err)
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			l.Printf("Error while closing stream with err=%s \n", err)
		}
	}(f)

	var decodedMap map[uuid.UUID]*models.Object
	d := json.NewDecoder(f)

	// Decoding the serialized data
	err = d.Decode(&decodedMap)
	if err != nil {
		l.Printf("Error while decoding storage with err=%s \n", err)
		return
	}

	//save to Vault
	models.ImportStorage(decodedMap)
}

// saveStorage exports map to a file
func saveStorage(l *log.Logger) {

	storage := models.ExportStorage()

	//encoding map to save
	b := new(bytes.Buffer)

	e := json.NewEncoder(b)

	// Encoding the map
	err := e.Encode(storage)
	if err != nil {
		l.Printf("Error while encoding storage with err=%s \n", err)
		return
	}

	//saving file
	f, err := os.Create("import/data.json")

	if err != nil {
		l.Printf("Error while creating file with err=%s \n", err)
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			l.Printf("Error while closing stream with err=%s \n", err)
		}
	}(f)

	//writing encoded data to stream
	_, err = f.WriteString(b.String())

	if err != nil {
		l.Printf("Error while writing to file with err=%s \n", err)
		return
	}

	l.Println("Successful file export")
}
