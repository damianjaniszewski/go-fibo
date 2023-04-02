package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func fibo(num int64) int64 {
	if num <= 0 {
		return 0
	} else if num == 1 {
		return 1
	}
	return fibo(num-1) + fibo(num-2)
}

func handlerSequence(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)

	// Stop here if its Preflighted OPTIONS request
	if r.Method == "GET" {

		sequence, err := strconv.ParseInt(vars["sequence"], 10, 64)
		if err != nil {
			log.Errorf("%s", err)
		}

		start := time.Now()
		res := fibo(sequence)
		elapsed := time.Since(start)
		log.WithFields(logrus.Fields{
			"sequence": sequence, "result": res, "elapsed": elapsed,
		}).Infof("fibo sequence: %d, result: %d, elapsed: %v", sequence, res, elapsed.Seconds())

		mComputeTime.With(prometheus.Labels{"guid": cfg.InstanceGUID, "sequence": vars["sequence"]}).Set(elapsed.Seconds())

		fmt.Fprintf(w, "{\"sequence\":%d, \"result\":%d, \"elapsed\":%v}", sequence, res, elapsed.Seconds())
	}
}

func handlerFibo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	w.WriteHeader(http.StatusOK)

	// Stop here if its Preflighted OPTIONS request
	if r.Method == "GET" {
		sequence := cfg.FiboMin + rand.Intn(cfg.FiboMax-cfg.FiboMin)
		log.WithFields(logrus.Fields{"sequence": sequence}).Debugf("generated sequence: %d", sequence)

		start := time.Now()
		res := fibo(int64(sequence))
		elapsed := time.Since(start)
		log.WithFields(logrus.Fields{
			"sequence": sequence, "result": res, "elapsed": elapsed,
		}).Infof("fibo sequence: %d, result: %d, elapsed: %v", sequence, res, elapsed.Seconds())

		mComputeTime.With(prometheus.Labels{"guid": cfg.InstanceGUID, "sequence": strconv.Itoa(sequence)}).Set(elapsed.Seconds())

		fmt.Fprintf(w, "{\"sequence\":%d, \"result\":%d, \"elapsed\":%v}", sequence, res, elapsed.Seconds())
	}
}
