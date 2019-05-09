package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// http://blog.golang.org/error-handling-and-go
type appHandler func(http.ResponseWriter, *http.Request) *appError

type Report struct {
	Cpu   string `json:"cpu"`
	Memory string `json:"memory"`
	Storage  string `json:"storage"`
    Deviceid  string `json:"deviceid"`
}

func ConfigureDB(ctx context.Context) (*bigquery.Client, error) {
	projectId := os.Getenv("PROJ_ID")
	if projectId == "" {
		return nil, errors.New("NO PROJECT ID DEFINED")
	}
	client, err := bigquery.NewClient(ctx, projectId)
	return client, err
}

func AddReport(r Report, client *bigquery.Client, ctx context.Context) error {
	datasetId := os.Getenv("DATASET_ID")
	tableId := os.Getenv("TABLE_ID")
	if datasetId == "" || tableId == "" {
		return errors.New("DATASET_ID or TABLE_ID is not defined")
	}
	u := client.Dataset(datasetId).Table(tableId).Uploader()

	err := u.Put(ctx, fr)
	if err != nil {
		if multiError, ok := err.(bigquery.PutMultiError); ok {
			for _, err1 := range multiError {
				for _, err2 := range err1.Errors {
					fmt.Println(err2)

				}
			}
			return multiError
		} else {
			fmt.Println(err)
		}
	}
	return nil
}

// This is doing streaming inserts which is a cost of 0.010 USD / 200 MBs for US mult-region
// See : https://cloud.google.com/bigquery/pricing for more details
func JsonToBigQueryStorer(w http.ResponseWriter, r *http.Request) *appError {
	ctx := context.Background()
	client, err := ConfigureDB(ctx)
	if err != nil {
		return appErrorf(err, "Firestore error : %v", err)
	}
	defer client.Close()
	decoder := json.NewDecoder(r.Body)
	var report Report
	err = decoder.Decode(&report)
	if err != nil {
		return appErrorf(err, "Cannot Decode Data : %v", err)
	}
	err = AddReport(report, client, ctx)
	if err != nil {
		return appErrorf(err, "Could not add report to database : %v", err)
	}
	fmt.Fprintf(w, "Success")
	return nil
}

type appError struct {
	Error   error
	Message string
	Code    int
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		log.Printf("Handler error: status code: %d, message: %s, underlying err: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func registerHandlers() {
	r := mux.NewRouter()

	r.Methods("POST").Path("/report").
		Handler(appHandler(JsonToBigQueryStorer))
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}
	registerHandlers()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}