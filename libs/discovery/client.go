package discovery

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// CheckIn registers the local server to the discovery server.
func CheckIn(url, name, ip string, port int) error {
	_, err := http.Post(url+fmt.Sprintf("/services/%s/ip", name), "text/plain", strings.NewReader(ip))
	if err != nil {
		return err
	}

	_, err = http.Post(url+fmt.Sprintf("/services/%s/port", name), "text/plain", strings.NewReader(fmt.Sprintf("%d", port)))
	if err != nil {
		return err
	}
	return nil
}

// Like CheckIn but it tries to connect every 2 seconds
func CheckInBlock(url, name, ip string, port int) error {
	for {
		err := CheckIn(url, name, ip, port)
		if err != nil {
			log.Printf("Error while checking in, sleeping 2s: %v", err)
			time.Sleep(2 * time.Second)
		} else {
			return nil
		}
	}
	return nil
}

func CheckOut(url, name string) error {
	req, err := http.NewRequest("DELETE", url+fmt.Sprintf("/services/%s", name), nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}
