package traefik

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type Client struct {
	Endpoint string
}

func (c *Client) PutConfiguration(configuration Configuration) error {
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
