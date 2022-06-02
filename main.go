package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KforG/p2pool-scanner-go/logging"
)

func main() {
	//init logging
	logFile, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.SetLogFile(logFile)

	logging.Infof("P2pool-Scannner-Go started up! \n")

	// Hold up program for Go routines and exit gracefully
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	logging.Infof("Input detected, exiting application..\n")
}
