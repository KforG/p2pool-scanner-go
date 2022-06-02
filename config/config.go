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
	Port           string   `json:"Port"`
	BootstrapNodes []string `json:"BootstrapNodes"`
	GeoLocation    struct {
		API        string `json:"API"`
		AcessKey   string `json:"AcessKey"`
		Parameters string `json:"Parameters"`
	} `json:"GeoLocation"`
	Domain struct {
		Check   bool
		Domains []string `json:"Domains"`
	} `json:"Domain"`
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

	if len(Active.Domain.Domains) > 0 {
		Active.Domain.Check = true
	} else {
		Active.Domain.Check = false
	}

	return nil
}
