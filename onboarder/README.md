#### Onboarder:
This webservice is used to register PnP clients.

Client details consists of `{MacId, OperationType}`

To run the webservice : `[go run onboarder.go --registry_address=<consul_registry_address> --server_name="ClientOnboardService" --server_address <consul_registry_address>:<some_port>]`

##### Sample Rest calls:
1. <I>Register Client : </I> POST : `[http://172.16.128.147:8099/pnp/clients?MacId=00:0c:29:c0:2b:a8&OpType=deploySatellite]` (This POST call can also be used to update client's op-type)
2. <I>Get all Client details: </I> GET : `[http://172.16.128.147:8099/pnp/clients]`
3. <I>Get Client detail: </I> GET : `[http://172.16.128.147:8099/pnp/clients/{mac}]`
4. <I>Deregister Client : </I> DELETE : `[http://172.16.128.147:8099/pnp/clients/{mac}]`
