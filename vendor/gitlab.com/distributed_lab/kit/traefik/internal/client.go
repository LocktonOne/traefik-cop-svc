package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
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

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		return nil
	case http.StatusNotFound, http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return errors.New("invalid request")
	default:
		return errors.New(
			"something bad happened")
	}
}
