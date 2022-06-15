package geo

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/KforG/p2pool-scanner-go/config"
	"github.com/KforG/p2pool-scanner-go/logging"
	"github.com/KforG/p2pool-scanner-go/util"
)

type Geo struct {
	Location struct {
		CountryFlag string `json:"country_flag"`
	} `json:"location"`
	CountryCode string `json:"country_code"`
	RegionName  string `json:"region_name"`
	CountryName string `json:"country_name"`
}

//This function queries a iplookup tool for the geolocation of the node
//It does not take in a URL but an IP address
func GetGeoLocation(nodeIPPort string, jsonPayload *Geo) error {
	s := strings.SplitN(nodeIPPort, ":", len(nodeIPPort))
	nodeIP := s[0]
	err := util.GetJson(fmt.Sprintf(config.Active.GeoLocation.API+nodeIP+config.Active.GeoLocation.AcessKey+config.Active.GeoLocation.Parameters), &jsonPayload)
	if err != nil {
		return err
	}
	return err
}

func GetFlag(countryCode string, flagURL string) (path string, err error) {
	//Check if the flag already exist in the subdirectory /flags/, if it does there's no need to fetch it again
	flag, path := FlagExist(countryCode)
	if flag {
		return path, nil
	}

	resp, err := http.Get(flagURL)
	if err != nil {
		logging.Errorf("Failed to fetch flag %s\n %s\n", flagURL, err)
		return "", err
	}
	defer resp.Body.Close()

	//Create a file to write the fetched data into
	wd, err := os.Getwd()
	if err != nil {
		logging.Errorf("Error creating flag, couldn't find working directory\n %s\n", err)
		return "", err
	}
	fp := "/flags/"
	if _, err = os.Stat(filepath.Join(wd, fp)); os.IsNotExist(err) {
		logging.Infof("No 'flags' subdirectory found, creating..\n")
		err = os.Mkdir(filepath.Join(wd, fp), 0700)
		if err != nil {
			logging.Errorf("Error couldn't create subdirectory! Launch the program with correct permissions or create /flags/ before launch!")
			return "", err
		}
	}

	file, err := os.Create(filepath.Join(wd, fmt.Sprintf("/flags/%s.svg", countryCode)))
	if err != nil {
		logging.Errorf("Error couldn't create file containing flag!\n %s\n", err)
		return "", err
	}
	defer file.Close()

	//Dumps the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		logging.Errorf("Failed to dump flag body into file, %s\n", err)
		return "", err
	}
	_, path = FlagExist(countryCode)
	return path, nil
}

//This function simply checks if the flag already exist, returns a bool and the path to the flag
func FlagExist(countryCode string) (exist bool, path string) {
	wd, err := os.Getwd()
	if err != nil {
		logging.Errorf("Error checking if flag exists, couldn't find working directory\n %s\n", err)
		return false, ""
	}

	//If another GeoLocater is used this function might need changing to the file name
	flagsPath := fmt.Sprintf("/flags/%s.svg", countryCode)
	path = filepath.Join(wd, flagsPath)
	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		return false, ""
	}
	return true, path
}
