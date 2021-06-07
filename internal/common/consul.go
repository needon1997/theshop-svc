package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/needon1997/theshop-svc/internal/common/config"
	"net/http"
)

const REGISTER_URI = "/v1/agent/service/register"
const DEREGISTER_URI = "/v1/agent/service/deregister/%s"
const SERVICES_FILTER_URI = "/v1/agent/services?filter=Service==\"%s\""

func RegisterSelfToConsul() error {
	conf := config.ServerConfig.ConsulConfig
	if conf.Check.CheckMethod != "HTTP" {

	}
	if conf.Host == "" {
		conf.Host = config.ServerConfig.Host
	}
	if conf.Port == 0 {
		conf.Port = config.ServerConfig.Port
	}
	var checkUrl string
	if conf.Check.CheckMethod == "HTTP" {
		checkUrl = fmt.Sprintf("http://%s:%v%s", conf.Host, conf.Port, conf.Check.Uri)
	}
	if conf.Check.CheckMethod == "GRPC" {
		checkUrl = fmt.Sprintf("%s:%v", conf.Host, conf.Port)
	}
	body := map[string]interface{}{
		"Name":    conf.Name,
		"Tags":    conf.Tags,
		"Address": conf.Host,
		"Port":    conf.Port,
		"ID":      conf.Id,
		"Check": map[string]interface{}{
			conf.Check.CheckMethod: checkUrl,
			"Interval":             conf.Check.Interval,
			"Method":               conf.Check.Method,
		},
	}
	bodyByte, err := json.Marshal(body)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("PUT", conf.Url+REGISTER_URI, bytes.NewBuffer(bodyByte))
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("server return %v", response.StatusCode))
	}
	return nil
}

type Service struct {
	Address string `json:"Address"`
	Port    int    `json:"Port"`
}

func GetServicesByNameTags(svcName string, tag string) ([]Service, error) {
	conf := config.ServerConfig.ConsulConfig
	url := conf.Url + fmt.Sprintf(SERVICES_FILTER_URI, svcName)
	if tag != "" {
		url += "&filter= \"%s\" in Tags"
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("server return %s", response.StatusCode))
	}
	svcsmap := make(map[string]Service)
	err = json.NewDecoder(response.Body).Decode(&svcsmap)
	if err != nil {
		return nil, err
	}
	svcs := make([]Service, len(svcsmap))
	for _, v := range svcsmap {
		svcs = append(svcs, v)
	}
	return svcs, nil
}
func DeRegisterFromConsul() error {
	conf := config.ServerConfig.ConsulConfig
	request, err := http.NewRequest("PUT", conf.Url+fmt.Sprintf(DEREGISTER_URI, conf.Id), nil)
	if err != nil {
		return err
	}
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("server return %s", response.StatusCode))
	}
	return nil
}
