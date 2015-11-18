package main

import (
    "os"
    "github.com/op/go-logging"
)

var log = logging.MustGetLogger("instaprinter")

var format = logging.MustStringFormatter(
    "%{color}[%{time:15:04:05.000}] (%{program}:%{shortfunc}:%{pid}) - %{message}%{color:reset}",
)

func setupLogging() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	logging.SetBackend(backendFormatter)
}
