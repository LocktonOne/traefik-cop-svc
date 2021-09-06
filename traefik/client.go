package traefik

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	redis "github.com/go-redis/redis/v8"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Client interface {
	PutConfiguration(configuration Configuration) error
}

type RestClient struct {
	Endpoint string
}

func (c *RestClient) PutConfiguration(configuration Configuration) error {
	fmt.Println("putting config by rest api")

	url := fmt.Sprintf("%s/api/providers/rest", c.Endpoint)
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(&configuration); err != nil {
		return errors.Wrap(err, "failed to marshal body")
	}

	req, err := http.NewRequest(http.MethodPut, url, &body)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to send request")
	}
	defer response.Body.Close()

	respBB, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	if response.StatusCode > 300 {
		return errors.From(errors.New("failed to update configuration"),
			logan.F{"status": response.StatusCode,
				"body": string(respBB),
			})
	}

	return nil
}

type RedisClient struct {
	ctx context.Context
	rc  *redis.Client
}

func (c *RedisClient) PutConfiguration(configuration Configuration) error {
	fmt.Println("putting config to redis")

	rawConfiguration := mergeDisjointMaps(craftRoutersMap(configuration), craftServersMap(configuration))
	return c.rc.MSet(c.ctx, rawConfiguration).Err()
}

func craftServersMap(c Configuration) map[string]string {
	servers := make(map[string]string)

	base := "traefik/http/services"
	for serviceName, service := range c.HTTP.Services {
		serversLoadBalancerServerBase := fmt.Sprintf("%s/%s/loadbalancer/servers", base, serviceName)
		for i, server := range service.LoadBalancer.Servers {
			serverUrlKey := fmt.Sprintf("%s/%d/url", serversLoadBalancerServerBase, i)
			servers[serverUrlKey] = server.URL
		}
	}

	return servers
}

func craftRoutersMap(c Configuration) map[string]string {
	routers := make(map[string]string)

	base := "traefik/http/routers"

	for routerName, router := range c.HTTP.Routers {
		routerBase := fmt.Sprintf("%s/%s", base, routerName)

		routers[fmt.Sprintf("%s/rule", routerBase)] = router.Rule
		routers[fmt.Sprintf("%s/service", routerBase)] = router.Service
		if router.Priority != 0 {
			routers[fmt.Sprintf("%s/priority", routerBase)] = fmt.Sprintf("%d", router.Priority)
		}
	}

	return routers
}

func mergeDisjointMaps(m1, m2 map[string]string) map[string]string {
	merged := m1

	for k, v := range m2 {
		merged[k] = v
	}

	return merged
}
