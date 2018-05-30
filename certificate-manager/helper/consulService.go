package helper

import (
	"os"
	"net/http"
	"fmt"
	log "github.com/ZTP/pnp/util/color"
	"io/ioutil"
	"encoding/json"
)

const consulPort = "8500"
var ConsulServiceName string

type ConsulCatalogService []struct {
	ServiceAddress string   `json:"ServiceAddress"`
	ServicePort    int  	`json:"ServicePort"`
}

func (c *ConsulCatalogService) GetServiceDetails() error {
	interfaceName := os.Getenv("SDP_NETWORK_INTERFACE")
	if interfaceName == "" {
		log.Fatalf("Provide \"SDP_NETWORK_INTERFACE\" environment variable")
	}
	httpClient := &http.Client{}
	onboarderUrl := fmt.Sprintf("http://%s:%s/v1/catalog/service/%s", GetIPFromIPwithCIDR(GetIPv4ForInterfaceName(interfaceName).String()), consulPort, ConsulServiceName)
	req, err := http.NewRequest("GET", onboarderUrl, nil)
	if err != nil {
		return err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	consulServiceDetail, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(consulServiceDetail, &c); err != nil {
		return err
	}
	return nil
}
