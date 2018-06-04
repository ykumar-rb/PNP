package handlers

import (
	"context"
	"io/ioutil"
	"github.com/golang/protobuf/ptypes"
	pb "github.com/ZTP/pnp/common/proto"
	proto "github.com/ZTP/certificate-manager/proto/certificate"
	pnpproto "github.com/ZTP/pnp/pnp-proto"
	"github.com/ZTP/pnp/common/color"
	"github.com/ZTP/certificate-manager/helper"
	"github.com/ZTP/onboarder/handlers"
	"github.com/go-redis/redis"
	"os"
	"net/http"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
)

type PnPCertificateService struct {}

type ClientConfig struct {
	MacId	string	`json:"MacId"`
	OpType	string	`json:"OpType"`
}

func (s *PnPCertificateService) GetCertificates (ctx context.Context, clientInfo *proto.ClientInfo, certificateResponse *proto.ServerCertificate) (err error) {
	var responseMsg string
	clientMAC := clientInfo.CommonClientInfo.ClientInfo.MACAddr
	color.Printf("Received certificate request from client, Client MAC: %v", clientMAC)
	responseMsg, clientEnvName := getClientDetails(clientMAC)
	pwd, _ := os.Getwd()
	certFile, err := ioutil.ReadFile(pwd+"/certs/server.crt")
	if err != nil {
		color.Fatalf("Error reading certificate file", err)
		return err
	}
	encryptCertFile := helper.Encrypt([]byte(certFile), clientMAC)
	*certificateResponse = proto.ServerCertificate{CommonServerResponse: &pnpproto.CommonServerResponse{
		ResponseHeader: &pb.ResponseHeader{Identifiers: &pb.Identifiers{TraceID: clientInfo.
			CommonClientInfo.RequestHeader.Identifiers.TraceID, MessageID: clientInfo.
			CommonClientInfo.RequestHeader.Identifiers.MessageID}, ResponseTimestamp:
				ptypes.TimestampNow()},}, ServerCert: encryptCertFile, ClientEnvName: clientEnvName,
					ResponseMessage: responseMsg}
	return
}

func getClientDetails (mac string) (string, string) {
	var onboarderUrl string
	httpClient := &http.Client{}
	clientDetail := ClientConfig{}
	consulCatalogService := helper.ConsulCatalogService{}
	clientEnv := NewEnv()

	err := consulCatalogService.GetServiceDetails()
	if err != nil {
		color.Warnf("Error", err)
		return "", ""
	}
	onboarderUrl = fmt.Sprintf("http://%s:%s/pnp/clients/%s", consulCatalogService[0].ServiceAddress, strconv.Itoa(consulCatalogService[0].ServicePort), mac)
	req, err := http.NewRequest("GET", onboarderUrl, nil)
	if err != nil {
		color.Warnf("Error", err)
		return "", ""
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		color.Warnf("Error", err)
		return "", ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		color.Warnf("Client \"%v\" not registered", mac)
		return "", ""
	}
	clientDetailByt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.Warnf("Error", err)
		return "", ""
	}
	if err := json.Unmarshal(clientDetailByt, &clientDetail); err != nil {
		color.Warnf("Error", err)
		return "", ""
	}
	if clientEnvVal := clientEnv.RedisClient.HGet(mac, "EnvName").String(); clientEnvVal != "" {
		clientEnvName := strings.Split(clientEnvVal, ": ")
		return "ok", clientEnvName[len(clientEnvName)-1]
	}
	return "", ""
}

func NewEnv () *handlers.InstallEnv {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		color.Fatalf("Provide \"REDIS_ADDR\" environment variable")
	}
	return &handlers.InstallEnv{
		RedisClient: redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: "", // no password set
			DB:       0,  // use default DB
		}),
	}
}

