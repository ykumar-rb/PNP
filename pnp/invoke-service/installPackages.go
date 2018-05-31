package invoke

import (
	"time"
	"io"
	"fmt"
	"log"
	"github.com/ZTP/pnp/common"
	"golang.org/x/net/context"
	"github.com/ZTP/pnp/executor"
	proto "github.com/ZTP/pnp/pnp-proto"
)

func setPkgMsgType(serverPkgOperType proto.ServerMsgType, exeErr error) (clientPkgMsgType proto.ClientMsgType) {
	switch serverPkgOperType {
	case proto.ServerMsgType_IS_PKG_INSTALLED:
		{
			if exeErr == nil {
				clientPkgMsgType = proto.ClientMsgType_PKG_INSTALLED
			} else {
				clientPkgMsgType = proto.ClientMsgType_PKG_NOT_INSTALLED
			}
		}
	case proto.ServerMsgType_INSTALL_PKG, proto.ServerMsgType_INSTALL_PKG_FROM_FILE:
		{
			if exeErr == nil {
				clientPkgMsgType = proto.ClientMsgType_PKG_INSTALL_SUCCESS
			} else {
				clientPkgMsgType = proto.ClientMsgType_PKG_INSTALL_FAILED
				fmt.Printf("\nFailed to install package\n")
			}
		}
	case proto.ServerMsgType_GET_NEXT_PKG:
		{
			clientPkgMsgType = proto.ClientMsgType_PKG_ZTP_INIT
		}
	}
	return
}

func InitPkgMgmt(pnpClient proto.PnPService, clientInfo proto.ClientInfo) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	stream, err := pnpClient.GetPackages(cxt)
	clientMsgType := proto.ClientMsgType_PKG_ZTP_INIT

	clientMsg := &proto.ClientPkgRequest{CommonClientInfo: &proto.CommonClientInfo{RequestHeader:
	common.NewReqHdrGenerateTraceAndMessageID(), ClientInfo: &clientInfo},
		ClientMsgType: clientMsgType}
	serverPkgResp := &proto.ServerPkgResponse{}

	for {
		if err = stream.Send(clientMsg); err != nil {
			log.Fatalf("Failed to send client message, Error: %v", err)
		}

		serverPkgResp, err = stream.Recv()
		if err == io.EOF {
			fmt.Println("\nClosing connection...")
			stream.Close()
			break
		}

		if err != nil {
			fmt.Printf("Error while receiving data from server %v\n",  err)
		}

		var exeErr error
		if serverPkgResp.CommonServerResponse.GetServerCmdType() == proto.ServerCmdType_RUN {
			cmdStr := serverPkgResp.ServerInstructionPayload.Cmd
			exeErr = executor.ExecuteServerInstructions(cmdStr)
		}

		clientMsgType = setPkgMsgType(serverPkgResp.GetServerMsgType(), exeErr)

		traceId := serverPkgResp.CommonServerResponse.ResponseHeader.Identifiers.TraceID

		clientMsg = &proto.ClientPkgRequest{CommonClientInfo: &proto.CommonClientInfo{RequestHeader:
		common.NewReqHdrGenerateMessageID(traceId), ClientInfo: &clientInfo},
			ClientMsgType: clientMsgType }
	}
}
