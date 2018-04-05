package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
}

func crashHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	if r.Method == "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("/crash request received, crashing...")
	}
}

func restHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	vars := mux.Vars(r)

	// Stop here if its Preflighted OPTIONS request
	if r.Method == "POST" {

		sequence, err := strconv.ParseInt(vars["sequence"], 10, 64)
		if err != nil {
			log.Println(err)
		}

		start := time.Now()
		res := fibo(sequence)
		elapsed := time.Since(start)
		log.Printf("fibo(%d)=%d [%v]\n", sequence, res, elapsed)

		fmt.Fprintf(w, "{\"sequence\":%d, \"result\":%d}", sequence, res)
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/{sequence:[0-9]+}", restHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/health", healthHandler).Methods("GET", "OPTIONS")
	router.HandleFunc("/crash", crashHandler).Methods("POST", "OPTIONS")

	port := ":" + os.Getenv("PORT")
	log.Println("go-fibo started, listening on", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func fibo(num int64) int64 {
	if num <= 0 {
		return 0
	} else if num == 1 {
		return 1
	}
	return fibo(num-1) + fibo(num-2)
}
