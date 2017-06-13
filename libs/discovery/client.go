package discovery

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Name string
	Ip   string
	Port int
}

// CheckIn registers the local server to the discovery server.
func (s *Service) CheckIn(url string) error {
	_, err := http.Post(url+fmt.Sprintf("/services/%s/ip", s.Name), "text/plain", strings.NewReader(s.Ip))
	if err != nil {
		return err
	}

	_, err = http.Post(url+fmt.Sprintf("/services/%s/port", s.Name), "text/plain", strings.NewReader(fmt.Sprintf("%d", s.Port)))
	if err != nil {
		return err
	}
	return nil
}

// Like CheckIn but it tries to connect every 2 seconds
func (s *Service) CheckInBlock(url string) error {
	for {
		err := s.CheckIn(url)
		if err != nil {
			log.Printf("Error while checking in, sleeping 2s: %v", err)
			time.Sleep(2 * time.Second)
		} else {
			return nil
		}
	}
}

func (s *Service) CheckOut(url string) error {
	req, err := http.NewRequest("DELETE", url+fmt.Sprintf("/services/%s", s.Name), nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

func NewService(name, ip string, port int) *Service {
	return &Service{Name: name, Ip: ip, Port: port}
}

func ListServices(url string) (map[string]Service, error) {
	services := make(map[string]Service)
	resp, err := http.Get(url + "/services")
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return services, nil
	}
	if err != nil {
		return services, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return services, nil
	}

	for _, line := range strings.Split(string(content), "\n") {
		if line == "" {
			break
		}
		resp, _ = http.Get(fmt.Sprintf("%s/services/%s/ip", url, line))
		data, _ := ioutil.ReadAll(resp.Body)
		ip := string(data)

		resp, _ = http.Get(fmt.Sprintf("%s/services/%s/port", url, line))
		data, _ = ioutil.ReadAll(resp.Body)
		port, _ := strconv.Atoi(string(data))

		services[line] = Service{Name: line, Ip: ip, Port: port}
	}
	return services, nil

}
