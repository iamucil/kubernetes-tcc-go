package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"module": "main"})

var (
	revisionID     = "unknown"
	buildTimestamp = "unknown"
)

func main() {
	_, err := fmt.Fprintf(os.Stdout, "Blog service revisionID %s, built at %s\n", revisionID, buildTimestamp)

	if err != nil {
		panic(err)
	}

	for _, flag := range os.Args[1:] {
		if flag == "--version" {
			return
		}
	}

	logLevelStr := os.Getenv("LOG_LEVEL")
	if logLevelStr != "" {
		logLevel, err := logrus.ParseLevel(logLevelStr)
		if err != nil {
			panic(err)
		}
		logrus.SetLevel(logLevel)
	}

	apiServicePort := os.Getenv("API_PORT")
	if apiServicePort == "" {
		apiServicePort = "8080"
	}

	var apiServiceBasePath string
	if baseStr, found := os.LookupEnv("BASE_PATH"); found {
		if baseStr != "" {
			apiServiceBasePath = baseStr
		}
	}

	restContainer := restful.DefaultContainer
	restContainer.EnableContentEncoding(true)

	http.HandleFunc("/"+apiServiceBasePath, func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Helo, World!"))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("OK"))
	})

	var shuttingDown bool
	shutdownSignal := make(chan os.Signal)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	apiServer := &http.Server{Addr: ":" + apiServicePort}

	go func() {
		log.Infof("Service is ready")
		err := apiServer.ListenAndServe()
		if err != nil && (err != http.ErrServerClosed || !shuttingDown) {
			log.Fatalf("REST API Server Error. %v", err)
		}
	}()

	<-shutdownSignal
	shuttingDown = true
	log.Infof("Shutting down the server")

	go func() {
		<-shutdownSignal
		os.Exit(0)
	}()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	apiServer.Shutdown(shutdownCtx)
	log.Infof("Done...")
}
