package common

import (
	"strings"
	"runtime"
	"log"
	"fmt"
	"github.com/ZTP/pnp/executor"
	proto "github.com/ZTP/pnp/pnp-proto"
)

func PopulateClientDetails() (clientInfo proto.ClientInfo) {
	archType := runtime.GOARCH
	osType := runtime.GOOS
	getOSFlavorCmd := "lsb_release -a | grep Description | awk -F':' '{print $2}'"

	osFlavor, err := executor.ExecuteCommand(getOSFlavorCmd)
	if err != nil {
		log.Fatalf("Error while getting OS type: %v", err)
	}
	// ToDo: Client ID generation...
	clientId := "client1"

	clientInfo = proto.ClientInfo{OsType: osType, ArchType: archType, OsFlavor: osFlavor, ClientId: clientId}
	return
}

func ExecuteServerInstructions(cmdString []string) (exeErr error) {
	var errStr string
	cmd := strings.Join(cmdString, " && ")
	errStr, exeErr = executor.ExecuteCommand(cmd)
	if exeErr != nil {
		fmt.Printf("\nCommand <%v> failed to execute\nErrorString: %v\nError: %v\n", cmd, errStr, exeErr)
	}
	return exeErr
}
