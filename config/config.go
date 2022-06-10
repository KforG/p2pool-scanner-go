package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/KforG/p2pool-scanner-go/logging"
)

var Active configStruct

type configStruct struct {
	Currency       string   `json:"Currency"`
	CurrencyCode   string   `json:"CurrencyCode"`
	BootstrapNodes []string `json:"BootstrapNodes"`
	Port           string   `json:"Port"`
	RescanTime     int      `json:"RescanTime"`
	GeoLocation    struct {
		API        string `json:"API"`
		AcessKey   string `json:"AcessKey"`
		Parameters string `json:"Parameters"`
	} `json:"GeoLocation"`
	KnownDomains struct {
		Check      bool
		NodeDomain []struct {
			IP         string
			DomainName string `json:"DomainName"`
		} `json:"NodeDomain"`
	} `json:"KnownDomains"`
	WebPort string `json:"WebPort"`
}

func ReadConfig() error {
	logging.Infof("Reading config file...\n")
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		if os.IsNotExist(err) {
			logging.Fatalf("No config.json found in directory! Please create config.json using example_config.json\n")
			return err
		}
		logging.Errorf("Error reading config.json file! %s\n", err)
		return err
	}

	err = json.Unmarshal(file, &Active)
	if err != nil {
		logging.Errorf("Error unmarshaling configuration file %s\n", err)
		return err
	}

	if len(Active.KnownDomains.NodeDomain) > 0 {
		Active.KnownDomains.Check = true
	} else {
		Active.KnownDomains.Check = false
	}

	return nil
}
