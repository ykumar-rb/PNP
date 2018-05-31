package invoke_service

import (
	"time"
	"golang.org/x/net/context"
	proto "github.com/ZTP/certificate-manager/proto/certificate"
	pnpproto "github.com/ZTP/pnp/pnp-proto"
	"github.com/ZTP/pnp/common"
	"github.com/ZTP/pnp/common/color"
	"github.com/ZTP/certificate-manager/helper"
)

func GetCertificate (pnpClient proto.CertificateService, clientInfo pnpproto.ClientInfo) []byte{
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	clientMsg := &proto.ClientInfo{CommonClientInfo: &pnpproto.CommonClientInfo{RequestHeader:
	common.NewReqHdrGenerateTraceAndMessageID(), ClientInfo: &clientInfo}}

	color.Printf("\n\n[ CLIENT: FETCH CERTIFICATE ] Sending Request Message With MAC : %v\n\n", clientInfo.MACAddr)
	serverCertificate, err := pnpClient.GetCertificates(ctx, clientMsg)
	if err != nil {
		color.Fatalf("Error while receiving server certificate, Error: ", err)
	}
	color.Printf("Response from server: %v", serverCertificate.ResponseMessage)
	if serverCertificate.ResponseMessage == ""{
		color.Fatalf("Server could not find details for this client")
	}
	certificateBytes := serverCertificate.ServerCert
	decryptCert := helper.Decrypt(certificateBytes, clientInfo.MACAddr)
	return decryptCert
}
