package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"gopkg.in/yaml.v2"
)

// ServiceDef place holder for reading data from yaml
type ServiceDef map[string]map[string]interface{}

// Service is placeholder for the grpc service
type Service struct {
	Name        string
	Port        int
	Tag         string
	URL         string
	ContextPath string
	HealthCheck string
}

const (
	port        = "port"
	url         = "url"
	tag         = "tag"
	contextPath = "contextPath"
	healthCheck = "healthCheck"
)

// ReadServerConfigFromYaml reads yml form uri or local
func ReadServerConfigFromYaml(yamlPath string) ([]Service, error) {
	log.Println(fmt.Sprintf("Loading config path: %v", yamlPath))
	var serviceDef ServiceDef
	var bodyBytes []byte
	var readErr error
	if isURL, err := regexp.MatchString(`https?:\/\/.*`, yamlPath); err == nil && isURL {
		resp, err := http.Get(yamlPath)
		log.Println(resp)
		if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			bodyBytes, readErr = ioutil.ReadAll(resp.Body)
		}
	} else {
		bodyBytes, readErr = ioutil.ReadFile(yamlPath)
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

		if value, exists := details[url]; exists {
			srv.URL = value.(string)
		}

		if value, exists := details[contextPath]; exists {
			srv.ContextPath = value.(string)
		}

		if value, exists := details[healthCheck]; exists {
			srv.HealthCheck = value.(string)
		}

		if srv.ContextPath != "" && srv.HealthCheck != "" {
			srv.HealthCheck = fmt.Sprintf("%s%s", srv.ContextPath, srv.HealthCheck)
		}

		lis = append(lis, *srv)
	}
	return lis, nil
}
