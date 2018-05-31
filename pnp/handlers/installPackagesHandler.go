package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"log"
	"sync"
	"encoding/gob"
	"github.com/golang/protobuf/ptypes"
	"github.com/ZTP/pnp/util/server"
	pb "github.com/ZTP/pnp/common/proto"
	proto "github.com/ZTP/pnp/pnp-proto"
	"github.com/ZTP/pnp/common/color"
)

type PnPService struct {}

type ClientEnv struct {
	ClientConfigFile string
	AutoUpdate bool
}

type InstallEnv struct {
	mux sync.Mutex
	ClientEnvMap map[string]ClientEnv
}

func setPkgServerResponse (pkg server.Package,
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
	packageInfo := &server.PackageInfo{}
	installEnv := InstallEnv{}

	initialClientMsg, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		fmt.Printf("Error reading data from client, Error : %v", err)
		return err
	}
	serverPkgResponse = &proto.ServerPkgResponse{CommonServerResponse: &proto.CommonServerResponse{ServerCmdType: proto.ServerCmdType_INFO}}
	if err = stream.Send(serverPkgResponse); err != nil {
		fmt.Printf("Error while sending response to client, Error: %v", err)
		return err
	}
	pwd,_ := os.Getwd()
	installEnv.deSerializeStruct(pwd+"/../clientEnvMap.gob")
	clientIntructionFile := installEnv.fetchClientInstructionFileName(initialClientMsg.CommonClientInfo.ClientInfo.MACAddr)
	log.Printf("Instruction file for client %v : %v ", initialClientMsg.CommonClientInfo.ClientInfo.MACAddr ,clientIntructionFile)
	if err = server.GetConfigFromToml(clientIntructionFile, packageInfo); err != nil {
		log.Fatalf("Unable to get client instruction data from JSON file, Error: %v", err)
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

func (i *InstallEnv) fetchClientInstructionFileName (clientMac string) string {
	if i.ClientEnvMap[clientMac].ClientConfigFile == ""{
		color.Warnf("No instruction file found for the client : %v", clientMac)
	}
	return i.ClientEnvMap[clientMac].ClientConfigFile
}

func (i *InstallEnv) deSerializeStruct(serializedFile string) error {
	file, err := os.Open(serializedFile)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(i)
	}
	file.Close()
	return err
}
