package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"log"
	"sync"
	"strconv"
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/ptypes"
	"github.com/ZTP/pnp/util/server"
	pb "github.com/ZTP/pnp/common/proto"
	proto "github.com/ZTP/pnp/pnp-proto"
	"github.com/ZTP/pnp/common/color"
	"strings"
)

type PnPService struct {}

type ClientEnv struct {
	ClientConfigFile string
	AutoUpdate bool
	EnvName string
}

type InstallEnv struct {
	mux            sync.Mutex
	RedisClient *redis.Client
	clientEnv ClientEnv
}

var isPkgUpgrade bool

func getNextPackage(numPkgsToInstall int) (cmdType proto.ServerCmdType,
	serverMsgType proto.ServerMsgType) {

	if numPkgsToInstall == 0 {
		fmt.Println("\nDone with all pkgs\n")
		cmdType = proto.ServerCmdType_CLOSE_CONN
	} else {
		cmdType = proto.ServerCmdType_INFO
		serverMsgType = proto.ServerMsgType_GET_NEXT_PKG
	}
	return
}

func setPkgServerResponse (pkg server.Package,
	clientMsgType proto.ClientMsgType, numPkgsToInstall int, autoUpdate bool) (cmdType proto.ServerCmdType,
	serverMsgType proto.ServerMsgType, exeCmd []string){

	switch clientMsgType {
	case proto.ClientMsgType_PKG_ENV_INIT:
		{
			cmdType = proto.ServerCmdType_RUN
			serverMsgType = proto.ServerMsgType_INITIALIZE_ENV
			exeCmd = pkg.ExportEnv
		}
	case proto.ClientMsgType_PKG_ENV_INITIALIZE_FAILED:
		{
			fmt.Printf("ENV initialize failed..")
			cmdType = proto.ServerCmdType_CLOSE_CONN
		}
	case proto.ClientMsgType_PKG_ENV_INITIALIZED:
		{
			cmdType = proto.ServerCmdType_RUN
			serverMsgType = proto.ServerMsgType_IS_PKG_INSTALLED
			exeCmd = pkg.CheckInstalledCmd
		}
	case proto.ClientMsgType_PKG_INSTALLED:
		{
			fmt.Printf("Package %v installed.. Checking if it is latest version\n", pkg.Name)
			cmdType = proto.ServerCmdType_RUN
			serverMsgType = proto.ServerMsgType_IS_PKG_OUTDATED
			exeCmd = pkg.IsPackageOutdated
		}
	case proto.ClientMsgType_PKG_VERSION_OUTDATED:
		{
			isPkgUpgrade = true
			fmt.Printf("Package %v installed is outdated..\n", pkg.Name)

			if autoUpdate {
				cmdType = proto.ServerCmdType_RUN
				exeCmd = pkg.UninstallPackage
			} else {
				cmdType = proto.ServerCmdType_MANUAL_UPDATE

				uninstStr := strings.Join(pkg.UninstallPackage, ",")
				instStr := strings.Join(pkg.InstallInstructions, ",")
				rollBckStr := strings.Join(pkg.RollbackPackage, ",")

				combStr := uninstStr + "#" + instStr + "#" + rollBckStr
				exeCmd[0] = combStr
			}
			serverMsgType = proto.ServerMsgType_UNINSTALL_PKG
		}
	case proto.ClientMsgType_PKG_UNINSTALL_FAILED:
		{
			fmt.Printf("Uninstallation of package %v failed\n", pkg.Name)
			cmdType = proto.ServerCmdType_CLOSE_CONN
		}
	case proto.ClientMsgType_PKG_UNINSTALL_SUCCESS:
		{
			fmt.Printf("Uninstall package %v success\n", pkg.Name)
			cmdType = proto.ServerCmdType_RUN

			if pkg.UpdateRepo != nil {
				exeCmd = pkg.UpdateRepo
			}

			for _, cmd := range pkg.InstallInstructions {
				exeCmd = append(exeCmd, cmd)
			}
			serverMsgType = proto.ServerMsgType_INSTALL_PKG
		}
	case proto.ClientMsgType_PKG_NOT_INSTALLED:
		{
			if autoUpdate {
				cmdType = proto.ServerCmdType_RUN
			} else {
				cmdType = proto.ServerCmdType_MANUAL_UPDATE
			}

			if pkg.UpdateRepo != nil {
				exeCmd = pkg.UpdateRepo
			}

			for _, cmd := range pkg.InstallInstructions {
				exeCmd = append(exeCmd, cmd)
			}
			serverMsgType = proto.ServerMsgType_INSTALL_PKG
		}
	case proto.ClientMsgType_PKG_VERSION_LATEST:
		{
			fmt.Printf("Package %v is latest..", pkg.Name)
			cmdType, serverMsgType = getNextPackage(numPkgsToInstall)
		}
	case proto.ClientMsgType_PKG_INSTALL_SUCCESS:
		{
			fmt.Printf("Package %v installed\n", pkg.Name)
			cmdType, serverMsgType = getNextPackage(numPkgsToInstall)
		}
	case proto.ClientMsgType_PKG_INSTALL_FAILED:
		{
			if isPkgUpgrade {
				cmdType = proto.ServerCmdType_RUN
				serverMsgType = proto.ServerMsgType_ROLLBACK_PKG
				exeCmd = pkg.RollbackPackage
			} else {
				fmt.Printf("Installation of package %v failed\n", pkg.Name)
				cmdType = proto.ServerCmdType_CLOSE_CONN
			}
		}
	case proto.ClientMsgType_PKG_ROLLBACK_SUCCESS:
		{
			fmt.Printf("Package %v rollback success\n", pkg.Name)
			cmdType, serverMsgType = getNextPackage(numPkgsToInstall)
		}
	case proto.ClientMsgType_PKG_ROLLBACK_FAILED:
		{
			fmt.Printf("Package %v rollback failed\n", pkg.Name)
			cmdType = proto.ServerCmdType_CLOSE_CONN
		}
	case proto.ClientMsgType_GET_NEXT:
		{
			serverMsgType = proto.ServerMsgType_GET_NEXT_PKG
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
		color.Warnf("Error reading data from client, Error : %v", err)
		return err
	}

	installEnv.RedisClient = initializeClient()
	clientIntructionFile, err := installEnv.fetchClientInstructionFileName(initialClientMsg.CommonClientInfo.ClientInfo.MACAddr)
	if err != nil {
		return err
	}
	log.Printf("Instruction file for client %v : %v ", initialClientMsg.CommonClientInfo.ClientInfo.MACAddr ,clientIntructionFile)
	if err = server.GetConfigFromToml(clientIntructionFile, packageInfo); err != nil {
		color.Warnf("Unable to get client instruction data from JSON file, Error: %v", err)
		return
	}

	serverPkgResponse = &proto.ServerPkgResponse{CommonServerResponse: &proto.CommonServerResponse{ResponseHeader:
	&pb.ResponseHeader{Identifiers: &pb.Identifiers{TraceID: initialClientMsg.CommonClientInfo.RequestHeader.Identifiers.TraceID,
		MessageID: initialClientMsg.CommonClientInfo.RequestHeader.Identifiers.MessageID}, ResponseTimestamp:
	ptypes.TimestampNow()}, ServerCmdType:
		proto.ServerCmdType_INFO}, ServerMsgType: proto.ServerMsgType_CLIENT_AUTHENTICATED}

	if err = stream.Send(serverPkgResponse); err != nil {
		fmt.Printf("Error while sending response to client, Error: %v", err)
		return err
	}

	numPkgsToInstall := len(packageInfo.Packages)

	for _, pkg := range packageInfo.Packages {
		numPkgsToInstall = numPkgsToInstall - 1
		isPkgUpgrade = false
		for {
			clientPkgMsg, err := stream.Recv()
			if err == io.EOF {
				return nil
			}

			if err != nil {
				fmt.Printf("Error reading data from client, Error : %v", err)
				goto label
			}
			cmdType, pkgOperType, exeCmd := setPkgServerResponse(pkg, clientPkgMsg.GetClientMsgType(), numPkgsToInstall,
				installEnv.clientEnv.AutoUpdate)

			serverPkgResponse = &proto.ServerPkgResponse{CommonServerResponse: &proto.CommonServerResponse{ResponseHeader:
				&pb.ResponseHeader{Identifiers: &pb.Identifiers{TraceID:
					clientPkgMsg.CommonClientInfo.RequestHeader.Identifiers.TraceID, MessageID:
						clientPkgMsg.CommonClientInfo.RequestHeader.Identifiers.MessageID}, ResponseTimestamp:
							ptypes.TimestampNow()}, ServerCmdType: cmdType}, ServerInstructionPayload:
								&proto.ServerInstructionPayload{exeCmd}, ServerMsgType: pkgOperType,
									PackageDetails: &proto.PackageDetails{PackageName: pkg.Name, PackageVersion: pkg.Version,
									AutoUpdate: installEnv.clientEnv.AutoUpdate}}

			if err = stream.Send(serverPkgResponse); err != nil {
				fmt.Printf("Error while sending response to client, Error: %v", err)
				goto label
			}

			if cmdType == proto.ServerCmdType_CLOSE_CONN {
				goto label
			}

			if pkgOperType == proto.ServerMsgType_GET_NEXT_PKG {
				break
			}
		}
	}
	label:
	stream.Close()
	return err
}

func (i *InstallEnv) fetchClientInstructionFileName (clientMac string) (string, error) {
	clientEnvName := i.RedisClient.HGet(clientMac, "EnvName").Val()
	i.clientEnv.EnvName = clientEnvName
	color.Printf("ENV name from mac: %v:%v", clientMac,clientEnvName)
	clientEnvAutoUpdate := i.RedisClient.HGet(clientMac, "AutoUpdate").Val()// string: true/false
	instructionFileName,err := i.RedisClient.Get(clientEnvName).Result()
	if err != nil {
		color.Warnf("Error while fetching Environment filename, Error : %v", err)
		return "", err
	}
	i.clientEnv.AutoUpdate,_ = strconv.ParseBool(clientEnvAutoUpdate)
	if err != nil {
		color.Warnf("Error while converting string to boolean, Error : %v", err)
		return "", err
	}
	i.clientEnv.ClientConfigFile = instructionFileName
	return instructionFileName, nil
}

func initializeClient () (*redis.Client) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		color.Fatalf("Provide \"REDIS_ADDR\" environment variable")
	}
	Client_EnvPath := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return Client_EnvPath
}
