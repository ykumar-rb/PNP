package invoke

import (
	"time"
	"io"
	"fmt"
	"log"
	"github.com/ZTP/pnp/common"
	"golang.org/x/net/context"
	"github.com/ZTP/pnp/executor"
	"github.com/go-redis/redis"
	proto "github.com/ZTP/pnp/pnp-proto"
	"strings"
)

type SoftwareDB struct {
	Name string
	Version string
	AvailVersion string
	Action string
	Status string
	Install string
	UnInstall string
	Rollback string
}

var DBClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6389",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func SetDataInDB(client *redis.Client, SDB *SoftwareDB) {
	client.HSet(SDB.Name, "Name", SDB.Name)
	client.HSet(SDB.Name, "Version", SDB.Version)
	client.HSet(SDB.Name, "AvailVersion", SDB.AvailVersion)
	client.HSet(SDB.Name, "Action", SDB.Action)
	client.HSet(SDB.Name, "Status", SDB.Status)
	client.HSet(SDB.Name, "Install", SDB.Install)
	client.HSet(SDB.Name, "UnInstall", SDB.UnInstall)
	client.HSet(SDB.Name, "Rollback", SDB.Rollback)
}

func setPkgMsgType(serverPkgResp proto.ServerPkgResponse, exeErr error) (clientPkgMsgType proto.ClientMsgType) {
	switch serverPkgResp.GetServerMsgType() {
	case proto.ServerMsgType_CLIENT_AUTHENTICATED:
		{
			clientPkgMsgType = proto.ClientMsgType_PKG_ENV_INIT
		}
	case proto.ServerMsgType_INITIALIZE_ENV:
		{
			if exeErr == nil {
				clientPkgMsgType = proto.ClientMsgType_PKG_ENV_INITIALIZED
			} else {
				clientPkgMsgType = proto.ClientMsgType_PKG_ENV_INITIALIZE_FAILED
			}
		}
	case proto.ServerMsgType_IS_PKG_INSTALLED:
		{
			if exeErr == nil {
				clientPkgMsgType = proto.ClientMsgType_PKG_INSTALLED
			} else {
				clientPkgMsgType = proto.ClientMsgType_PKG_NOT_INSTALLED
			}
		}
	case proto.ServerMsgType_IS_PKG_OUTDATED:
		{
			if exeErr == nil {
				clientPkgMsgType = proto.ClientMsgType_PKG_VERSION_LATEST
			} else {
				clientPkgMsgType = proto.ClientMsgType_PKG_VERSION_OUTDATED
			}
		}
	case proto.ServerMsgType_UNINSTALL_PKG:
		{
			if serverPkgResp.CommonServerResponse.ServerCmdType != proto.ServerCmdType_MANUAL_UPDATE {
				if exeErr == nil {
					clientPkgMsgType = proto.ClientMsgType_PKG_UNINSTALL_SUCCESS
				} else {
					clientPkgMsgType = proto.ClientMsgType_PKG_UNINSTALL_FAILED
					SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
					serverPkgResp.PackageDetails.PackageVersion, AvailVersion: serverPkgResp.PackageDetails.PackageVersion,
						Action: "NOACTION", Status:"Uninstall while upgrade Failed", Install: "-", UnInstall: "-", Rollback: "-"})
				}
			} else {
				commandArray := strings.Split(serverPkgResp.ServerInstructionPayload.Cmd[0], "#")
				version := DBClient.HGet(serverPkgResp.PackageDetails.PackageName, "Version").Val()
				SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
				version, AvailVersion: serverPkgResp.PackageDetails.PackageVersion, Action: "UPGRADE", Status: "-",
				Install: commandArray[1], UnInstall:
						commandArray[0], Rollback: commandArray[2]})

				clientPkgMsgType = proto.ClientMsgType_GET_NEXT
			}
		}
	case proto.ServerMsgType_INSTALL_PKG:
		{
			if serverPkgResp.CommonServerResponse.ServerCmdType != proto.ServerCmdType_MANUAL_UPDATE {
				if exeErr == nil {
					clientPkgMsgType = proto.ClientMsgType_PKG_INSTALL_SUCCESS
					SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
					serverPkgResp.PackageDetails.PackageVersion, AvailVersion: serverPkgResp.PackageDetails.PackageVersion,
						Action: "NOACTION", Status:"Install package Success", Install: "-", UnInstall: "-", Rollback: "-"})
				} else {
					clientPkgMsgType = proto.ClientMsgType_PKG_INSTALL_FAILED
					SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
					serverPkgResp.PackageDetails.PackageVersion, AvailVersion: serverPkgResp.PackageDetails.PackageVersion,
						Action: "NOACTION", Status:"Install Package Failed", Install: "-", UnInstall: "-", Rollback: "-"})
					fmt.Printf("\nFailed to install package\n")
				}
			} else {
				clientPkgMsgType = proto.ClientMsgType_GET_NEXT
				commandArray := strings.Join(serverPkgResp.ServerInstructionPayload.Cmd, ",")
				SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
				serverPkgResp.PackageDetails.PackageVersion, AvailVersion: serverPkgResp.PackageDetails.PackageVersion,
					Action: "INSTALL", Status:"-", Install: commandArray, UnInstall: "-", Rollback: "-"})
			}

		}
	case proto.ServerMsgType_ROLLBACK_PKG:
		{
			if exeErr == nil {
				clientPkgMsgType = proto.ClientMsgType_PKG_ROLLBACK_SUCCESS
				version := DBClient.HGet(serverPkgResp.PackageDetails.PackageName, "Version").Val()
				SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
				version, AvailVersion: "-", Action: "NOACTION", Status:"Rollback Success", Install: "-",
				UnInstall: "-", Rollback: "-"})
			} else {
				clientPkgMsgType = proto.ClientMsgType_PKG_ROLLBACK_FAILED
				version := DBClient.HGet(serverPkgResp.PackageDetails.PackageName, "Version").Val()
				SetDataInDB(DBClient, &SoftwareDB{Name: serverPkgResp.PackageDetails.PackageName, Version:
				version, AvailVersion: "-", Action: "NOACTION", Status:"Rollback Failed", Install: "-",
					UnInstall: "-", Rollback: "-"})
			}
		}
	case proto.ServerMsgType_GET_NEXT_PKG:
		{
			clientPkgMsgType = proto.ClientMsgType_PKG_ENV_INIT
		}
	}
	return
}

func InitPkgMgmt(pnpClient proto.PnPService, clientInfo proto.ClientInfo) {
	cxt, cancel := context.WithTimeout(context.Background(), time.Minute*20)
	defer cancel()
	stream, err := pnpClient.GetPackages(cxt)
	clientMsgType := proto.ClientMsgType_AUTHENTICATE_CLIENT
	var cmdStr []string

	clientMsg := &proto.ClientPkgRequest{CommonClientInfo: &proto.CommonClientInfo{RequestHeader:
	common.NewReqHdrGenerateTraceAndMessageID(), ClientInfo: &clientInfo},
		ClientMsgType: clientMsgType}
	serverPkgResp := &proto.ServerPkgResponse{}

	for {
		if err = stream.Send(clientMsg); err != nil {
			log.Fatalf("Failed to send client message, Error: %v", err)
		}

		serverPkgResp, err = stream.Recv()
		if err == io.EOF || serverPkgResp.CommonServerResponse.GetServerCmdType() == proto.ServerCmdType_CLOSE_CONN {
			fmt.Println("\nClosing connection...")
			stream.Close()
			break
		}

		if err != nil {
			fmt.Printf("Error while receiving data from server %v\n",  err)
		}

		var exeErr error

		if serverPkgResp.CommonServerResponse.GetServerCmdType() == proto.ServerCmdType_RUN {
			cmdStr = serverPkgResp.ServerInstructionPayload.Cmd
			if serverPkgResp.PackageDetails.AutoUpdate {
				exeErr = executor.ExecuteServerInstructions(cmdStr)
			}
		}

		clientMsgType = setPkgMsgType(*serverPkgResp, exeErr)

		traceId := serverPkgResp.CommonServerResponse.ResponseHeader.Identifiers.TraceID

		clientMsg = &proto.ClientPkgRequest{CommonClientInfo: &proto.CommonClientInfo{RequestHeader:
		common.NewReqHdrGenerateMessageID(traceId), ClientInfo: &clientInfo},
			ClientMsgType: clientMsgType }
	}
}
