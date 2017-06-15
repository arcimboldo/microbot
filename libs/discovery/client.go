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
	Name      string
	Ip        string
	Port      int
	done      chan int
	Url       string
	checkedIn bool
}

// CheckIn registers the local server to the discovery server.

func (s *Service) checkin() error {
	_, err := http.Post(s.Url+fmt.Sprintf("/services/%s", s.Name), "text/plain", strings.NewReader(fmt.Sprintf("%s:%d", s.Ip, s.Port)))
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) checkout() error {
	if !s.checkedIn {
		return fmt.Errorf("Service not checked in, check in first")
	}

	req, err := http.NewRequest("DELETE", s.Url+fmt.Sprintf("/services/%s", s.Name), nil)
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	return err
}

func (s *Service) CheckIn(url string) error {
	if s.checkedIn {
		return fmt.Errorf("Service already checked in at %q: check out first.", s.Url)
	}
	s.Url = url

	err := s.checkin()
	if err != nil {
		return err
	}

	s.done = make(chan int)

	go func() {
		for {
			select {
			case <-s.done:
				log.Println("Checking out as expected")
				s.CheckOut()
			case <-time.After(5 * time.Second):
				log.Println("Checking in again")
				s.checkin()
			}
		}
	}()

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

func (s *Service) CheckOut() error {
	s.done <- 1
	return nil
}

func NewService(name, ip string, port int) *Service {
	return &Service{Name: name, Ip: ip, Port: port, checkedIn: false}
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
		resp, err := http.Get(fmt.Sprintf("%s/services/%s/ip", url, line))
		if err != nil {
			log.Printf("Unable to get IP for service %q: %q", line, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Unable to get IP for service %q: not found", line)
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Unable to get IP for service %q: %q", line, err)
			continue
		}
		ip := string(data)

		resp, err = http.Get(fmt.Sprintf("%s/services/%s/port", url, line))
		if err != nil {
			log.Printf("Unable to get PORT for service %q: %q", line, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Unable to get PORT for service %q: not found", line)
			continue
		}
		data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Unable to get PORT for service %q: %q", line, err)
			continue
		}
		port, err := strconv.Atoi(string(data))
		if err != nil {
			log.Printf("Unable to get PORT for service %q: %q", line, err)
			continue
		}

		services[line] = Service{Name: line, Ip: ip, Port: port}
	}
	return services, nil

}
