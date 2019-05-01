package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/damianjaniszewski/zpages"
	// "zpages"

	"github.com/codingconcepts/env"
	"github.com/joho/godotenv"

	"github.com/gofrs/uuid"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	appName    = "go-fibo"
	appVersion = "v0.0.14"
)

var (
	cfg = config{}

	listenAddress string
	localTime     time.Time
	tzOffset      int

	logger = logrus.New()
	log    = logger.WithFields(logrus.Fields{})
)

// metrics
var (
	mServerInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_fibo_info",
			Help: "Information about go-fibo environment.",
		},
		[]string{"name", "version", "guid", "address", "start_time"})

	mComputeTime = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "go_fibo_compute_time_seconds",
			Help: "Fibo compute time.",
		},
		[]string{"guid", "sequence"})
)

func main() {

	// update mServerInfo periodically
	go func() {
		for {
			time.Sleep(cfg.MetricsUpdateTime)

			log.WithFields(logrus.Fields{
				"startTime": localTime.Format(time.RFC3339), "uptime": time.Since(localTime).Seconds(),
			}).Debugf("uptime: %v", time.Since(localTime).Seconds())

			mServerInfo.With(prometheus.Labels{
				"name": cfg.ApplicationName, "version": cfg.ApplicationVersion, "guid": cfg.InstanceGUID,
				"address": cfg.InstanceAddress, "start_time": localTime.Format(time.RFC3339),
			}).Set(time.Since(localTime).Seconds())
		}
	}()

	// http server gracefull shutdown
	chShutdown := make(chan os.Signal, 1)

	// zPages
	z := &zpages.Handler{
		Version:  zpages.Version{Module: cfg.ApplicationName, Version: cfg.ApplicationVersion},
		LogLevel: zpages.LogLevel{Log: cfg.LogLevel, Debug: cfg.DebugLevel, Format: cfg.LogAs},

		ServiceStatus: zpages.ServiceStatus{
			GUID: cfg.InstanceGUID, Name: cfg.ApplicationName, Type: "application", URI: cfg.InstanceAddress,
			Health: zpages.HealthStatus{Status: zpages.HealthStatusOK}, Readiness: zpages.ReadinessStatus{Status: zpages.ReadinessStatusReady},
			Updated: time.Now(),
		},

		ShutdownChannel: chShutdown,

		Logger: logger,
	}
	z.Init(cfg.ApplicationName, cfg.ApplicationVersion, cfg.InstanceGUID, cfg.InstanceAddress, zpages.ServiceTypeApplication,
		cfg.LogLevel, cfg.DebugLevel, cfg.LogAs)

	// http router
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/api/v1/sequence/{sequence:[0-9]+}", handlerSequence).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/v1/fibo", handlerFibo).Methods("GET", "OPTIONS")

	router.HandleFunc("/healthz", z.Healthz).Methods("GET", "OPTIONS")
	router.HandleFunc("/readyz", z.Readyz).Methods("GET", "OPTIONS")

	router.HandleFunc("/support/v1/quiesce", z.SupportQuiesce).Methods("POST", "OPTIONS")
	router.HandleFunc("/support/v1/resume", z.SupportResume).Methods("POST", "OPTIONS")
	router.HandleFunc("/support/v1/quit", z.SupportQuit).Methods("POST", "OPTIONS")
	router.HandleFunc("/support/v1/restart", z.SupportRestart).Methods("POST", "OPTIONS")
	router.HandleFunc("/support/v1/fail", z.SupportFail).Methods("POST", "OPTIONS")
	router.HandleFunc("/support/v1/crash", z.SupportCrash).Methods("POST", "OPTIONS")

	// router.HandleFunc("/support/v1/info", z.SupportInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/support/v1/version", z.SupportVersion).Methods("GET", "OPTIONS")
	router.HandleFunc("/support/v1/env", z.SupportEnv).Methods("GET", "OPTIONS")
	router.HandleFunc("/support/v1/loglevel", z.SupportLogLevel).Methods("GET", "PUT", "OPTIONS")
	// router.HandleFunc("/support/v1/icon", z.SupportIcon).Methods("GET", "OPTIONS")

	router.Handle("/metrics", promhttp.Handler())
	router.Use(middlewareLogging)

	// http server settings
	server := &http.Server{
		Handler: router,
		Addr:    listenAddress,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 20 * time.Second,
		ReadTimeout:  20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// listen without blocking
	go func() {
		log.WithFields(logrus.Fields{"listenAddress": listenAddress}).Infof("listening on %s", listenAddress)

		log.Errorf("%v", server.ListenAndServe())
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(chShutdown, os.Interrupt)

	// Block until we receive our signal.
	<-chShutdown

	log.Info("received SIGINT, shutting down")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownWaitTime)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	server.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.

	logger.Exit(0)
}

// log requests when logLevel set to TRACE
func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.WithFields(logrus.Fields{"request": *r}).Trace()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// initialization
func init() {
	if err := godotenv.Load(); err != nil {
		logger.Fatal("error loading .env file")
	}

	if err := env.Set(&cfg); err != nil {
		logger.Fatalf("%+v", err)
	}

	localTime = time.Now()

	if tzLocation, err := time.LoadLocation(cfg.Timezone); err == nil {
		localTime = localTime.In(tzLocation)
	}
	cfg.Timezone, tzOffset = localTime.Zone()

	if _, instGUIDFound := os.LookupEnv("CF_INSTANCE_GUID"); !instGUIDFound {
		cfg.InstanceGUID = uuid.Must(uuid.NewV4()).String()
	}

	listenAddress = ":" + cfg.InstancePort

	if _, instAddrFound := os.LookupEnv("CF_INSTANCE_ADDR"); !instAddrFound {
		cfg.InstanceAddress = cfg.Hostname + ":" + cfg.InstancePort
	}

	switch cfg.LogLevel {
	case "PANIC":
		logger.Level = logrus.PanicLevel
	case "FATAL":
		logger.Level = logrus.FatalLevel
	case "ERROR":
		logger.Level = logrus.ErrorLevel
	case "WARNING":
		logger.Level = logrus.WarnLevel
	case "INFO":
		logger.Level = logrus.InfoLevel
	case "DEBUG":
		logger.Level = logrus.DebugLevel
	case "TRACE":
		logger.Level = logrus.TraceLevel
	default:
		logger.Level = logrus.InfoLevel
	}

	logFullTimestamp := false
	if cfg.DebugLevel > 0 {
		logFullTimestamp = true
	}

	switch cfg.LogAs {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: logFullTimestamp})
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: logFullTimestamp})
	}
	logger.Out = os.Stdout

	if cfg.LogDetails {
		log = logger.WithFields(logrus.Fields{"app": cfg.ApplicationName, "version": cfg.ApplicationVersion, "instance": cfg.InstanceAddress, "guid": cfg.InstanceGUID})
	}

	log.WithFields(logrus.Fields{
		"min": cfg.FiboMin, "max": cfg.FiboMax, "logLevel": logger.Level, "debugLevel": cfg.DebugLevel, "logAs": cfg.LogAs,
	}).Infof("%s %s initialized", appName, appVersion)

	log.WithFields(logrus.Fields{
		"tzName": cfg.Timezone, "tzOffset": tzOffset,
	}).Debugf("timezone: %s, time offset: %d", cfg.Timezone, tzOffset)

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(mServerInfo)
	prometheus.MustRegister(mComputeTime)

	mServerInfo.With(prometheus.Labels{
		"name": cfg.ApplicationName, "version": cfg.ApplicationVersion, "guid": cfg.InstanceGUID,
		"address": cfg.InstanceAddress, "start_time": localTime.Format(time.RFC3339),
	}).Set(time.Since(localTime).Seconds())
}
