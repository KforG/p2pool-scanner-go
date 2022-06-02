package util

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KforG/p2pool-scanner-go/logging"
)

func GetJson(url string, target interface{}) error {
	resp, err := http.Get(fmt.Sprintf("http://%s", url))
	if err != nil {
		logging.Errorf("Error fetching data from %s\n", url)
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(target)
	if err != nil {
		logging.Errorf("Error parsing response from %s\n", url)
		return err
	}
	return err
}

func RemoveStringSliceIndex(i int, s *[]string) {
	if i != len(*s)-1 {
		(*s)[i] = (*s)[len(*s)-1]
	}
	// drop the last element
	*s = (*s)[:len(*s)-1]
}
