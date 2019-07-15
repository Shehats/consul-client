package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"gopkg.in/yaml.v2"
)

// ServiceDef place holder for reading data from yaml
type ServiceDef map[string]map[string]interface{}

// Service is placeholder for the grpc service
type Service struct {
	Name string
	Port int
	Tag  string
	URL  string
}

const (
	port = "port"
	url  = "url"
	tag  = "tag"
)

// ReadServerConfigFromYaml reads yml form uri or local
func ReadServerConfigFromYaml(path string) ([]Service, error) {
	fmt.Println(path)
	var serviceDef ServiceDef
	var bodyBytes []byte
	var readErr error
	if isURL, err := regexp.MatchString(`https?:\/\/.*`, path); err != nil && isURL {
		resp, err := http.Get(path)
		if err != nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			bodyBytes, readErr = ioutil.ReadAll(resp.Body)
		}
	} else {
		bodyBytes, readErr = ioutil.ReadFile(path)
	}

	if readErr == nil {
		yaml.Unmarshal(bodyBytes, &serviceDef)
	} else {
		return nil, readErr
	}

	lis := make([]Service, 0)
	for name, details := range serviceDef {
		fmt.Println(name)
		srv := new(Service)
		srv.Name = name
		if value, exists := details[port]; exists {
			srv.Port = value.(int)
		}
		if value, exists := details[tag]; exists {
			srv.Tag = value.(string)
		}
		if value, exists := details[name]; exists {
			srv.URL = value.(string)
		}
		lis = append(lis, *srv)
	}
	return lis, nil
}
