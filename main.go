package main

import (
	"crypto/tls"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	skipVerify = kingpin.Flag("skip-verify", "Will ignore certificate warnings").Bool()
	port       = kingpin.Flag("port", "Port to listen on").Default("8000").String()
	jamsHost   = kingpin.Flag("host", "Host for JAMS with scheme and port. e.g. http://localhost:6371").Required().String()
	pwd        = kingpin.Flag("password", "JAMS user password").Short('p').Required().String()
	username   = kingpin.Flag("username", "JAMS user username").Short('u').Required().String()
	level      = kingpin.Flag("log-level", "log level for application").Short('l').HintOptions("warn", "error", "info", "debug").Default("error").String()
)

func main() {
	kingpin.Version("v0.0.1")
	kingpin.Parse()
	switch *level {
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.ErrorLevel)
	}

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		handleProm(w, r)
	})
	logrus.Info("Starting server on :8000")
	logrus.Fatal("Fatal server error", http.ListenAndServe(":"+*port, nil))

}

func handleProm(w http.ResponseWriter, r *http.Request) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *skipVerify},
		},
	}
	jamsClient := JAMSClient{Username: *username, Password: *pwd, Client: &client, Host: *jamsHost}
	if err := jamsClient.Login(); err != nil {
		logrus.Fatal("Fatal error with Jams client", err)
	}
	registry := prometheus.NewRegistry()
	registerJAMSMetrics(registry, jamsClient)
	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
