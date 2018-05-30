# sdp-ztp

Plug-n-Play serves the following purpose:
1. Automated deployment of Service Delivery platform on a node in a given role.
2. Automated Software (e.g.: SDP core services) Update proposal with the option to enable automatic update on per node basis.
3. Should be able to perform custom node configuration.

#### SDP-ZTP consists of the following components:
1. ZTP server :- Automated deployment & configuration of components viz. Matchbox & Dnsmasq.
2. PnP Server & Client :- Components to assist installation of Service Delivery Platform on Master/Satellite nodes, perform update proposals.
3. Onboarder :- Webservice component to register, update, fetch & deregister PnP clients. Components such as PnP Server & Certificate Manager uses this webservice to fetch client details. An administrator can use this service to perform CRUD operations for a pnp client.
4. Certificate Manager :- Component to assist certificate management on the PnP-Server. This component provides PnP-client with the certificates to start a secure communication with PnP-server using Grpc & TLS.

#### Steps to setup SDP-ZTP:
1. To setup ZTP with matchbox, PnP & Certificate Manager components, follow the instructions present in (https://github.com/RiverbedTechnology/sdp-ztp/tree/master/ZTP).
2. To run only the PnP-Server/Client, Onboarder and Certificate Manager, run the Onboarder webservice first (https://github.com/RiverbedTechnology/sdp-ztp/tree/master/onboarder)  followed by CertificateManager service (https://github.com/RiverbedTechnology/sdp-ztp/tree/master/certificate-manager), then start the PnP Server & Client (https://github.com/RiverbedTechnology/sdp-ztp/tree/master/pnp).
