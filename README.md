It is a service that keeps a config for traefic in the RAM (map type of the `map[service_name]service_config`)
and updates the routing table of the traefic by sending this map directly to the traefic.
Each service announces in its config section `cop` in which it prescribes certain values,
this section becomes `service_config` in the already mentioned map. 
The name of the service turns into the key (`service_name`) for this service in the config map.
The `cop` section in the service config looks like a trace. way:

```yaml
cop:
    disabled: false
    endpoint: http://cop # endpoint of traefik_cop service
    upstream: http://marketplace # where to route traffic to requests
    service_name: the name of the service for that map in the cop and, accordingly, under this name the traffic will know this service
    service_port: 80
    service_prefix: /integrations/marketplace # this is the most interesting. 
    #According to this prefix, the cop will form a rule according to which traffic will determine that it needs to route the request to this service. 
    #In case service_prefix is ​​passed, the traffic will route all requests that came to the endpoint /integrations/marketplace/... 
    #(for another service, obviously, a different prefix) to http://marketplace (that is, the service’s upstream address, by which in the cluster/docker composite you can reach the service)
```

Prefixes are not transmitted in the horizon and api because all the connectors
go to their endpoints without any prefixes and therefore the cop explicitly compiles for them a set of rules for routing requests, which it then sends to traefic.