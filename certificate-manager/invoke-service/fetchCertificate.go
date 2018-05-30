package invoke_service

import (
	"time"
	"golang.org/x/net/context"
	proto "github.com/ZTP/certificate-manager/proto/certificate"
	pnpproto "github.com/ZTP/pnp/pnp-proto"
	"github.com/ZTP/pnp/common"
	"github.com/ZTP/pnp/util/color"
	"github.com/ZTP/certificate-manager/helper"
)

func GetCertificate (pnpClient proto.CertificateService, intf string) []byte{
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	MACAddr := helper.GetMACForInterfaceName(intf)
	defer cancel()
	clientInfo := common.PopulateClientDetails()
	clientMsg := &proto.ClientMACInfo{CommonClientInfo: &pnpproto.CommonClientInfo{RequestHeader:
	common.NewReqHdrGenerateTraceAndMessageID(), ClientInfo: &clientInfo}, MAC: MACAddr}
	//serverGetCertificateResponse := &proto.ServerCertificate{}

	color.Printf("\n\n[ CLIENT: FETCH CERTIFICATE ] Sending Request Message With MAC : %v\n\n", MACAddr)
	serverGetCertificateResponse, err := pnpClient.GetCertificates(ctx, clientMsg)
	if err != nil {
		color.Fatalf("Error while writing receiving server certificate ", err)
	}
	certificateBytes := serverGetCertificateResponse.ServerCert
	decrCert := helper.Decrypt(certificateBytes, MACAddr)
	/*ioutil.WriteFile("server.crt", certificateBytes, 0644)
	if err != nil {
		color.Fatalf("Error while writing certificate to file ", err)
	}*/
	return decrCert
}
