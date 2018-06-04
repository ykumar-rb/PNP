#!/bin/bash

set -e
GOLANG_VERSION=1.9
export PNP_USER=$1
export SDP_NETWORK_INTERFACE=$2
export PNP_USER_HOME="/home/$PNP_USER"
export PNP_USER_PROFILE="$PNP_USER_HOME/.profile"
export PNP_USER_GOPATH="$PNP_USER_HOME/go"
curr_dir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
onboarderRestApiPort=8099
redisPort=9999

setupPNPServer() {
    echo "Setting up PNP server ..."
    go get "github.com/BurntSushi/toml"
    pushd $PNP_USER_GOPATH/src/github.com/ZTP/pnp
    echo "Generating pnp-client binary..."
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o client client.go
    cp client /var/lib/matchbox/assets/coreos/client/    #Triggered from Preseed.cfg on client
    IP="$(ifconfig $SDP_NETWORK_INTERFACE | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')"
    echo "InterfaceName: $SDP_NETWORK_INTERFACE"
    echo "Starting PNP server..."
    export REDIS_ADDR=$IP:$redisPort
    go run server.go --registry_address=$IP --server_name "NewPnPService"  --cert_file "../certificate-manager/certs/server.crt" --key_file "../certificate-manager/certs/server.key" >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/pnpServer.log
    echo "PNP server setup done"
    popd
}

setupZTP() {
    #Matchbox setup
    wget -c https://github.com/coreos/matchbox/releases/download/v0.7.0/matchbox-v0.7.0-linux-amd64.tar.gz
    tar xzvf matchbox-v0.7.0-linux-amd64.tar.gz
    cp matchbox-v0.7.0-linux-amd64/matchbox /usr/local/bin
    mkdir -p /var/lib/matchbox/assets/coreos/ubuntu
    mkdir -p /var/lib/matchbox/groups/ubuntu
    mkdir -p /var/lib/matchbox/profiles
    mkdir -p /var/lib/matchbox/ignition
    mkdir -p /var/lib/matchbox/assets/coreos/client
    rm -rf matchbox-v0.7.0-linux-amd64*
    echo "Downloading initrd.gz"
    wget -cP /var/lib/matchbox/assets/coreos/client http://archive.ubuntu.com/ubuntu/dists/xenial-updates/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/initrd.gz
    echo "Downloading linux"
    wget -cP /var/lib/matchbox/assets/coreos/client http://archive.ubuntu.com/ubuntu/dists/xenial-updates/main/installer-amd64/current/images/netboot/ubuntu-installer/amd64/linux
    #Copy base resolv.conf
    cp /etc/resolv.conf /var/lib/matchbox/assets/coreos/client/
    #Dnsmasq setup
    mkdir -p /var/lib/tftpboot
    echo "Downloading undionly.kpxe"
    wget -cP /var/lib/tftpboot/ http://boot.ipxe.org/undionly.kpxe
    apt-get -y update && apt-get install -y dnsmasq
    configure_ZTP_services
}

setupCertificateManager() {
    echo "Setting up Certificate Manager..."
    echo "Fetching go libraries"
    go get "github.com/golang/protobuf/proto"
    go get "github.com/micro/go-micro"
    go get "github.com/micro/go-grpc"
    pushd $PNP_USER_GOPATH/src/github.com/ZTP/certificate-manager
    echo "Generating certificates..."
    go run GenerateTLSCertificate.go $SDP_NETWORK_INTERFACE
    cp certs/server.crt /var/lib/matchbox/assets/coreos/client/
    echo "Starting the certificate-manager service..."
    IP="$(ifconfig $SDP_NETWORK_INTERFACE | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')"
    go run server.go --registry_address=$IP --server_name="CertificateManagerService" --onboarder_service_name="ClientOnboardService" >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/certificateManager.log 2>&1 &
    popd
}

configure_ZTP_services() {
    mkdir -p $PNP_USER_GOPATH/src/github.com
    cp -r ${curr_dir}/../../ZTP $PNP_USER_GOPATH/src/github.com
    pushd $PNP_USER_GOPATH/src/github.com/ZTP/ztp-setup/sdp-ztp
    go run main.go >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/ztp.log 2>&1 &
    popd
}

setupConsul() {
    echo "Setting up consul"
    apt-get install -y zip
    wget -c https://releases.hashicorp.com/consul/1.0.7/consul_1.0.7_linux_amd64.zip
    unzip consul_1.0.7_linux_amd64.zip
    rm consul_1.0.7_linux_amd64.zip
    IP="$(ifconfig $SDP_NETWORK_INTERFACE | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')"
    ./consul agent -dev -bind=$IP -client $IP -ui -data-dir=/tmp/consul >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/consul.log 2>&1 &
    echo "Consul server running"
}

setupClientOnboarder() {
    echo "Setting up client-onboarder Rest Api"
    echo "Fetching go libraries"
    go get "github.com/micro/go-web"
    go get "github.com/emicklei/go-restful"
    go get "github.com/dotcloud/go-redis-server"
    go get "github.com/go-redis/redis"
    pushd $PNP_USER_GOPATH/src/github.com/ZTP/onboarder
    IP="$(ifconfig $SDP_NETWORK_INTERFACE | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')"
    echo "Starting Redis server"
    export REDIS_ADDR=$IP:$redisPort
    #BUG-FIX !! - Add "hvalues:  make(HashHash)," in github.com/dotcloud/go-redis-server/defaultHandler.go
    go run redisServer.go  >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/redisServer.log 2>&1 &
    go run onboarder.go --registry_address=$IP --server_name="ClientOnboardService" --server_address $IP:$onboarderRestApiPort >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/onboarder.log 2>&1 &
    popd
}

setupLoggingPath() {
    echo "Logging path being setup"
    mkdir -p $PNP_USER_GOPATH/src/github.com/ZTP/logs
    pushd $PNP_USER_GOPATH/src/github.com/ZTP/logs
    touch ztp.log consul.log onboarder.log certificateManager.log pnpServer.log serverUI.log redisServer.log
    popd
}

setupZTPServerUI() {
    echo "Setting up ZTP server UI"
    echo "Fetching go libraries..."
    go get "github.com/sirupsen/logrus"
    go get "github.com/spf13/viper"
    go get "github.com/tylerb/graceful"
    go get "github.com/carbocation/interpose"
    go get "github.com/gorilla/sessions"
    go get "github.com/jmoiron/sqlx"
    go get "golang.org/x/crypto/ssh/terminal"
    go get "github.com/gorilla/mux"
    cp -r ${curr_dir}/../../CloudInstall $PNP_USER_GOPATH/src/github.com
    pushd $PNP_USER_GOPATH/src/github.com/CloudInstall
    IP="$(ifconfig $SDP_NETWORK_INTERFACE | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}')"
    export ONBOARDER_API_ADDRESS=$IP:$onboarderRestApiPort
    go run server_ui.go >> $PNP_USER_GOPATH/src/github.com/ZTP/logs/serverUI.log 2>&1 &
    popd
}

is_go_installed() {
  [ ! -z "$(which go)" ]
}

is_curl_installed() {
    [ ! -z "$(which curl)" ]
}

install_curl() {
    apt-get -y update && apt-get -f install && apt-get -y install curl
}

install_go() {
  echo "Fetching go..."
  mkdir -p "$PNP_USER_GOPATH"

  pushd $(mktemp -d)
    curl -fL -o go.tgz "https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz"
    tar -C . -xzf go.tgz;
    mkdir -p /usr/lib/go-$GOLANG_VERSION
    mv go/* /usr/lib/go-$GOLANG_VERSION
  popd
    ln -s /usr/lib/go-$GOLANG_VERSION /usr/lib/go
    ln -s /usr/lib/go/bin/* /usr/bin/.
}

post_install() {
  echo "go installed..."
  echo "$(go version)"
  update_go_path
}

update_go_path() {
  if ! grep -q GOPATH $PNP_USER_PROFILE; then
    echo "export GOPATH=\"$PNP_USER_GOPATH\"" >> "$PNP_USER_PROFILE"
    echo 'export PATH="$PATH:$GOPATH/bin"' >> "$PNP_USER_PROFILE"
  fi
}

is_git_installed() {
    [ ! -z "$(which git)" ]
}

setupGit() {
  if is_git_installed; then
    echo "A version of git is already installed"
    echo "$(git version)"
  else
    apt-get install -y git
  fi
}

setupCurl() {
  if is_curl_installed; then
    echo "A version of curl is already installed"
    echo "$(curl --version)"
  else
    install_curl
  fi
}

only_run_as_root() {
    if [ "$(id -u)" -ne 0 ]; then
        echo "Error: must run as privileged user"
        exit 1
    fi
}

setupGo() {
  setupGit
  setupCurl
  if is_go_installed; then
    echo "A version of go is already installed"
    echo "$(go version)"
  else
    install_go
  fi
  post_install
}

if [ $# -ne 2 ]
then
    echo "Supply two arguments USER, Interface_name: e.g. \${USER} ens33"
    exit 1
fi

only_run_as_root
setupGo
setupLoggingPath
setupZTP
setupConsul
setupClientOnboarder
setupCertificateManager
setupZTPServerUI
setupPNPServer
