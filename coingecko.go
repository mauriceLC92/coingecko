package coingecko

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// DefaultEndpoint is the default API endpoint for the Coingecko public API
	DefaultEndpoint = "https://api.coingecko.com/api/v3"
)

type Client struct {
	HTTPClient *http.Client
	Endpoint   string
}

// Do performs the *http.Request and decodes the http.Response.Body into v and return the *http.Response. If v is an io.Writer it will copy the body to the writer.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	req.RequestURI = ""

	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}

	if c.Endpoint == "" {
		c.Endpoint = DefaultEndpoint
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		_, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return res, fmt.Errorf("http error code: %d", res.StatusCode)
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, res.Body)
		} else {
			decErr := json.NewDecoder(res.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				return nil, decErr
			}
		}
	}
	return res, nil
}

func (c *Client) Ping() (PingResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.Endpoint, "ping"), nil)
	if err != nil {
		return PingResponse{}, fmt.Errorf("new ping request: %w", err)
	}

	var pingResponse PingResponse
	_, err = c.Do(req, &pingResponse)
	if err != nil {
		return PingResponse{}, nil
	}

	return pingResponse, nil
}
