package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KforG/p2pool-scanner-go/config"
	"github.com/KforG/p2pool-scanner-go/handler"
	"github.com/KforG/p2pool-scanner-go/logging"
	"github.com/KforG/p2pool-scanner-go/scanner"
)

func main() {
	//init logging
	logFile, _ := os.OpenFile("debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logging.SetLogFile(logFile)

	logging.Infof("P2pool-Scannner-Go started up! \n")

	//init config
	err := config.ReadConfig()
	if err != nil {
		panic(err) //We can't continue the program if the config can't be loaded
	}
	logging.Infof("Config loaded successfully\n")

	// Start scanning
	nodes := scanner.Nodes{}
	go scanner.Scanner(&nodes)

	// API endpoint(s)
	go handler.Router(&nodes)

	// Hold up program for Go routines and exit gracefully
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	logging.Infof("Input detected, exiting application..\n")
}
