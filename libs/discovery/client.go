package discovery

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	Name      string
	Ip        string
	Port      int
	Url       string
	Rxs       []*regexp.Regexp
	done      chan int
	checkedIn bool
}

// CheckIn registers the local server to the discovery server.

func (s *Service) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Ip, s.Port)
}

func (s *Service) checkin() error {
	_, err := http.Post(fmt.Sprintf("%s/services/%s", s.Url, s.Name), "text/plain", strings.NewReader(fmt.Sprintf("%s:%d", s.Ip, s.Port)))
	if err != nil {
		return err
	}

	for i, re := range s.Rxs {
		url := fmt.Sprintf("%s/services/%s/re/%d", s.Url, s.Name, i)
		_, err = http.Post(url, "text/plain", strings.NewReader(re.String()))
		if err != nil {
			log.Printf("Ignoring error: unable to push  regexp %q to %q: %v", re, url, err)
			continue
		}
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

func (s *Service) AddRegexp(re string) error {
	rex, err := regexp.Compile(re)
	if err != nil {
		return err
	}
	s.Rxs = append(s.Rxs, rex)
	return nil
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

func ListServices(url string) (map[string]*Service, error) {
	services := make(map[string]*Service)
	resp, err := http.Get(url + "/services")
	if resp != nil && resp.StatusCode == http.StatusNotFound {
		return services, fmt.Errorf("No services found")
	}
	if err != nil {
		return services, err
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return services, nil
	}

	for _, name := range strings.Split(string(content), "\n") {
		if name == "" {
			break
		}

		resp, err := http.Get(fmt.Sprintf("%s/services/%s", url, name))
		if err != nil {
			log.Printf("Unable to get info for service %q: %q", name, err)
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("Unable to get IP for service %q: not found", name)
			continue
		}
		data, err := ioutil.ReadAll(resp.Body)

		servicedata := strings.Split(string(data), ":")
		if len(servicedata) != 2 {
			log.Printf("Error parsing service for %q: %q", name, string(data))
			continue
		}
		port, err := strconv.Atoi(servicedata[1])
		if err != nil {
			log.Printf("Wrong port for service %q: %q", name, servicedata[1])
			continue
		}
		newservice := Service{Name: name, Ip: servicedata[0], Port: port}
		log.Printf("Adding service %s at %s", newservice.Name, newservice.URL())
		services[name] = &newservice

		// Find regexps
		resp, err = http.Get(fmt.Sprintf("%s/services/%s/re", url, name))
		if err != nil {
			log.Printf("Unable to list regexps for service %q: %v", name, err)
			return services, nil
		}
		data, err = ioutil.ReadAll(resp.Body)
		indexes := []int{}
		for _, line := range strings.Split(string(data), "\n") {
			i, err := strconv.Atoi(line)
			if err == nil {
				indexes = append(indexes, i)
			}
		}

		sort.Ints(indexes)
		for _, i := range indexes {
			log.Printf("Regexp number %d for service %s", i, newservice.Name)
			resp, err := http.Get(fmt.Sprintf("%s/services/%s/re/%d", url, name, i))
			if err != nil {
				log.Printf("Unable to get regexp %q for service %q: %v", i, name, err)
				continue
			}
			data, _ := ioutil.ReadAll(resp.Body)
			line := string(data)
			re, err := regexp.Compile(line)
			if err != nil {
				log.Printf("Invalid regexp %q for service %q: %v", line, name, err)
			} else {
				newservice.Rxs = append(newservice.Rxs, re)
			}
		}
	}
	return services, nil

}
