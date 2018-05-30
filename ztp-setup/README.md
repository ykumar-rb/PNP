## ZTP & PnP setup

#### Prerequisites

1. Copy ZTP and PnP folders to home directory. e.g. `"/home/master"`
2. Run the following commands on :

    2.2. `cd /ZTP/`
    
    2.3. `chmod 755 ztp_setup.sh`
    
    2.4. `sudo ./ztp_setup.sh ${USER} <Interface_name>` e.g. `sudo ./ztp_setup.sh ${USER} ens33`
    
This setup configures ZTP on a new Ubuntu VM and then starts the PnP server. It involves the following steps:
1. GO installation
2. Install and configure matchbox & dnsmasq
3. Starts consul agent
4. Starts the PnP server

#### Post Installation
1. Create new custom virtual machine (PXE client) and Before booting it, get it's MAC address and register it using Onboarder Webservice. (Look Onboarder module Readme for more details on how to register)
2. Boot up the new PXE client.


When ever a new PXE client is booted, ZTP-setup installs ubuntu on the VM and then starts the PnP-client.
Once the PnP client is up and running, it automatically communicates with the PnP server for further instructions.(See PnP project for mode details)
