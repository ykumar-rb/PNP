package handlers

import(
	"os"
	"sync"
	"fmt"
	"io/ioutil"
	"github.com/emicklei/go-restful"
	"github.com/go-redis/redis"
	"encoding/json"
	log "github.com/ZTP/pnp/common/color"
	"net/http"
	"errors"
	"strconv"
	"github.com/ZTP/onboarder/config"
	"github.com/ZTP/onboarder/helper"
)

const clientList = "RegisteredClientList.toml"
type Onboarder struct{
	mux sync.Mutex
	clientList []byte
	clientListFile *os.File
}

type ClientConfig struct {
	MacId	string	`json:"MacId"`
	OpType	string	`json:"OpType"`
}

type ClientInfoList struct {
	ClientConfigs []ClientConfig
}

type InstallEnv struct {
	RedisClient    *redis.Client
	mux            sync.Mutex
}

func (o *Onboarder) GetAllRegisteredClients(req *restful.Request, rsp *restful.Response) {
	log.Printf("GET request : /pnp/clients")
	o.initFile()
	clientConfigList, err := o.getConfigs()
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	} else {
		rsp.WriteHeader(http.StatusOK)
		rsp.WriteEntity(clientConfigList)
	}
}

func (o *Onboarder) GetRegisteredClientDetails(req *restful.Request, rsp *restful.Response) {
	log.Printf("GET request : /pnp/clients/{mac}")
	var isClientReg bool
	o.initFile()
	clientInfoList, err := o.getConfigs()
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
		log.Fatalf("Error: ", err)
	}
	for i := 0; i < len(clientInfoList.ClientConfigs); i++ {
		if clientInfoList.ClientConfigs[i].MacId == req.PathParameter("mac") {
			isClientReg = true
			rsp.WriteHeader(http.StatusOK)
			rsp.WriteEntity(clientInfoList.ClientConfigs[i])
			break
		}
	}
	if isClientReg == false {
		rsp.WriteHeader(http.StatusNoContent)
	}
}

func (o *Onboarder) RegisterClient(req *restful.Request, rsp *restful.Response) {
	log.Printf("POST request : /pnp/clients?MacId&OpType")
	clientConfig := &ClientConfig{}
	o.initFile()
	clientConfig.MacId = req.QueryParameter("MacId")
	clientConfig.OpType = req.QueryParameter("OpType")
	err := o.addConfigToFile(*clientConfig)
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
	} else {
		rsp.WriteHeader(http.StatusCreated)
		rsp.WriteEntity("Client registered successfully")
	}
}

func (o *Onboarder) DeregisterClient(req *restful.Request, rsp *restful.Response) {
	log.Printf("DELETE request : /pnp/clients?MacId&OpType")
	clientInfoList := ClientInfoList{}
	o.initFile()
	clientInfoList, err := o.getConfigs()
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
		log.Fatalf("Error: ", err)
	}
	err = clientInfoList.deletefromClientList(req.PathParameter("mac"))
	if err != nil {
		rsp.WriteError(http.StatusNoContent, err)
	} else {
		o.mux.Lock()
		clientListJson, err := json.Marshal(clientInfoList)
		if err != nil {
			rsp.WriteError(http.StatusInternalServerError, err)
			log.Fatalf("Error: ", err)
		}
		o.clientListFile.Truncate(0)//empty prev contents of file.
		o.clientListFile.Seek(0,0)
		o.clientListFile.Write(clientListJson)
		o.clientListFile.Sync()
		defer o.mux.Unlock()
		defer o.clientListFile.Close()

		rsp.WriteHeader(http.StatusOK)
		rsp.WriteEntity("Client deregistered successfully")
	}
}

func (o *Onboarder) addConfigToFile(config ClientConfig) error {
	log.Printf("Writing to file ..")
	clientInfoList, err := o.getConfigs()
	if err != nil {
		return err
	}
	clientInfoList.addOrUpdateClientList(config)
	clientListJson, err := json.Marshal(clientInfoList)
	if err != nil {
		return err
	}
	o.mux.Lock()
	log.Printf("Filename : %v", o.clientListFile.Name())
	o.clientListFile.Truncate(0)//empty prev contents of file.
	o.clientListFile.Seek(0,0)
	o.clientListFile.Write(clientListJson)
	o.clientListFile.Sync()
	defer o.mux.Unlock()
	defer o.clientListFile.Close()
	return nil
}

func (o *Onboarder) getConfigs() (ClientInfoList, error) {
	log.Printf("Reading from file ..")
	clientInfoList := &ClientInfoList{}
	o.mux.Lock()
	if err := json.Unmarshal(o.clientList, &clientInfoList); err != nil {
		return ClientInfoList{}, err
	}
	defer o.mux.Unlock()
	return *clientInfoList, nil
}

func (l *ClientInfoList) addOrUpdateClientList (thisClientConfig ClientConfig) {
	var isClientReg bool
	for i := 0; i < len(l.ClientConfigs); i++ {
		conf := &l.ClientConfigs[i]
		if conf.MacId == thisClientConfig.MacId {
			conf.OpType = thisClientConfig.OpType
			isClientReg = true
			break
		}
	}
	if isClientReg == false {
		l.ClientConfigs = append(l.ClientConfigs, thisClientConfig)
	}
}

func (l *ClientInfoList) deletefromClientList (mac string) error {
	var isClientReg bool
	for i := 0; i < len(l.ClientConfigs); i++ {
		if l.ClientConfigs[i].MacId == mac {
			l.ClientConfigs = append(l.ClientConfigs[:i], l.ClientConfigs[i+1:]...)
			isClientReg = true
			break
		}
	}
	if isClientReg == false {
		return errors.New("Client was not registered")
	}
	return nil
}

func (o *Onboarder) initFile() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error: ", err)
	}
	clientConfFileName := pwd+"/"+clientList
	o.mux.Lock()
	clientListFile, err := os.OpenFile(clientConfFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error while opening file: ", err)
	}
	clientListBytes, err := ioutil.ReadAll(clientListFile)
	if err != nil {
		log.Fatalf("Error while converting file to bytes: ", err)
	}
	if len(clientListBytes) < 1 {
		clientListBytes = []byte("{}")
	}
	o.clientList = clientListBytes
	o.clientListFile = clientListFile
	defer o.mux.Unlock()
}

func (e *InstallEnv) CreateEnvironment (req *restful.Request, rsp *restful.Response) {
	log.Printf("POST request : /pnp/environment")
	newConfigEnv := &config.ConfigEnvironment{}
	requestByt, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
		log.Fatalf("", err)
	}
	if err := json.Unmarshal(requestByt, &newConfigEnv); err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
		log.Fatalf("", err)
	}
	log.Printf("Request contents: %v", newConfigEnv)
	err = e.StoreEnvPath(newConfigEnv)
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
		log.Fatalf("", err)
	}
	err = e.StoreMacEnv(newConfigEnv)
	if err != nil {
		rsp.WriteError(http.StatusInternalServerError, err)
		log.Fatalf("", err)
	}

	//defer e.mux.Unlock()
	rsp.WriteHeader(http.StatusOK)
}

func RegisterClient (mac string) error {
	httpClient := &http.Client{}
	consulCatalogService := helper.ConsulCatalogService{}
	err := consulCatalogService.GetServiceDetails()
	if err != nil {
		log.Warnf("Error", err)
		return err
	}
	onboarderUrl := fmt.Sprintf("http://%s:%s/pnp/clients/?MacId=%s", consulCatalogService[0].ServiceAddress, strconv.Itoa(consulCatalogService[0].ServicePort), mac)
	req, err := http.NewRequest("POST", onboarderUrl, nil)
	if err != nil {
		log.Warnf("Error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Warnf("Error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		log.Warnf("Client \"%v\" could not be registered", mac)
		return err
	}
	return nil
}

func (e *InstallEnv) UpdateEnvironment (req *restful.Request, rsp *restful.Response) {
	e.CreateEnvironment(req, rsp)
}

func (e *InstallEnv) StoreEnvPath(newEnv *config.ConfigEnvironment) error {
	e.mux.Lock()
	err := e.RedisClient.Set(newEnv.EnvironmentName, newEnv.InstructionFileName, 0).Err()
	if err != nil {
		return err
	}
	e.mux.Unlock()
	return nil
}

func (e *InstallEnv) StoreMacEnv(newEnv *config.ConfigEnvironment) error {
	e.mux.Lock()
	for i := 0; i < len(newEnv.Mac); i++ {
		e.RedisClient.HSet(newEnv.Mac[i],"EnvName",newEnv.EnvironmentName)
		e.RedisClient.HSet(newEnv.Mac[i],"AutoUpdate",newEnv.AutoUpdate)
		err := RegisterClient(newEnv.Mac[i])
		if err != nil {
			return err
		}
	}
	e.mux.Unlock()
	return nil
}
