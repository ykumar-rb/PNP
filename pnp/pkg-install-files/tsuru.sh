#!/usr/bin/env bash

TSURU_NOW_SCRIPT_URL="https://raw.githubusercontent.com/tsuru/now/master/run.bash"

# E.g.: "https://raw.githubusercontent.com/tsuru/tsuru/master/misc/git-hooks/pre-receive"
TSURU_NOW_HOOK_URL="https://raw.githubusercontent.com/tsuru/tsuru/master/misc/git-hooks/pre-receive"

# E.g.: "" or "--tsuru-from-source"
TSURU_NOW_OPTIONS="--tsuru-from-source"

curl -sL ${TSURU_NOW_SCRIPT_URL} > /tmp/tsuru-now.bash
chmod +x /tmp/tsuru-now.bash
sudo  \
  /tmp/tsuru-now.bash \
    --tsuru-pkg-${TSURU_MODE} \
    --hook-url ${TSURU_NOW_HOOK_URL} \
    --hook-name pre-receive \
    ${TSURU_NOW_OPTIONS}

if [ -d /usr/local/go ]; then
    export GOPATH=~/go
    mkdir -p $GOPATH
    if [ -f ~/.bashrc ]; then
        if ! grep 'export GOPATH' ~/.bashrc; then
          echo "Adding GOPATH=$GOPATH to ~/.bashrc"
          echo -e "export GOPATH=$GOPATH" | tee -a ~/.bashrc > /dev/null
        fi
    fi
fi

wget -q --no-check-certificate https://github.com/tsuru/tsuru-client/releases/download/1.5.1/tsuru_1.5.1_linux_386.tar.gz
tar -xzvf tsuru_1.5.1_linux_386.tar.gz
chmod +x tsuru
mv tsuru /usr/local/bin/