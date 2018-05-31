package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"log"
	"github.com/golang/protobuf/ptypes"
	"github.com/ZTP/pnp/common"
	"github.com/ZTP/pnp/config"
	pb "github.com/ZTP/pnp/common/proto"
	proto "github.com/ZTP/pnp/pnp-proto"
)

type PnPService struct {}

func setPkgServerResponse (pkg common.Package,
	clientMsgType proto.ClientMsgType, numPkgsToInstall int) (cmdType proto.ServerCmdType,
	serverMsgType proto.ServerMsgType, exeCmd []string){

	switch clientMsgType {
	case proto.ClientMsgType_PKG_ZTP_INIT:
		{
			cmdType = proto.ServerCmdType_RUN
			serverMsgType = proto.ServerMsgType_IS_PKG_INSTALLED
			exeCmd = pkg.CheckInstalledCmd
		}
	case proto.ClientMsgType_PKG_NOT_INSTALLED:
		{
			cmdType = proto.ServerCmdType_RUN
			if pkg.InstallFromFile != "" {
				serverMsgType = proto.ServerMsgType_INSTALL_PKG_FROM_FILE

			} else {
				if pkg.UpdateRepo != nil {
					exeCmd = pkg.UpdateRepo
				}
				serverMsgType = proto.ServerMsgType_INSTALL_PKG
				for _, cmd := range pkg.InstallInstructions {
					exeCmd = append(exeCmd, cmd)
				}
			}
		}
	case proto.ClientMsgType_PKG_INSTALLED:
		{
			fmt.Printf("Package %v already installed\n", pkg.Name)
			if numPkgsToInstall == 0 {
				cmdType = proto.ServerCmdType_CLOSE_CONN
			} else {
				cmdType = proto.ServerCmdType_INFO
				serverMsgType = proto.ServerMsgType_GET_NEXT_PKG
			}
		}
	case proto.ClientMsgType_PKG_INSTALL_SUCCESS:
		{
			fmt.Printf("Package %v installed\n", pkg.Name)
			if numPkgsToInstall == 0 {
				fmt.Println("\nDone with all pkgs\n")
				cmdType = proto.ServerCmdType_CLOSE_CONN
			} else {
				cmdType = proto.ServerCmdType_INFO
				serverMsgType = proto.ServerMsgType_GET_NEXT_PKG
			}
		}
	case proto.ClientMsgType_PKG_INSTALL_FAILED:
		{
			fmt.Printf("Installation of package %v failed\n", pkg.Name)
			cmdType = proto.ServerCmdType_CLOSE_CONN
		}
	}
	return
}

func (s *PnPService) GetPackages (ctx context.Context, stream proto.PnP_GetPackagesStream) (err error) {
	serverPkgResponse := &proto.ServerPkgResponse{}
	packageInfo := &common.PackageInfo{}
	pwd, _ := os.Getwd()

	if err = common.GetConfigFromToml(pwd + config.PackageFilePath, packageInfo); err != nil {
		log.Fatalf("Unable to get config data from JSON file, Error: %v", err)
	}

	numPkgsToInstall := len(packageInfo.Packages)

	for _, pkg := range packageInfo.Packages {
		numPkgsToInstall = numPkgsToInstall - 1

		for {
			clientPkgMsg, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				fmt.Printf("Error reading data from client, Error : %v", err)
				break
			}

			cmdType, pkgOperType, exeCmd := setPkgServerResponse(pkg, clientPkgMsg.GetClientMsgType(), numPkgsToInstall)

			serverPkgResponse = &proto.ServerPkgResponse{CommonServerResponse: &proto.CommonServerResponse{ResponseHeader:
			&pb.ResponseHeader{Identifiers: &pb.Identifiers{TraceID: clientPkgMsg.CommonClientInfo.RequestHeader.Identifiers.TraceID,
				MessageID: clientPkgMsg.CommonClientInfo.RequestHeader.Identifiers.MessageID}, ResponseTimestamp:
			ptypes.TimestampNow()}, ServerCmdType: cmdType}, ServerInstructionPayload:
			&proto.ServerInstructionPayload{exeCmd},
				ServerMsgType: pkgOperType}

			if err = stream.Send(serverPkgResponse); err != nil {
				fmt.Printf("Error while sending response to client, Error: %v", err)
				break
			}

			if pkgOperType == proto.ServerMsgType_GET_NEXT_PKG {
				break
			}
		}
		if err != nil {
			break
		}
	}
	stream.Close()
	return nil
}
