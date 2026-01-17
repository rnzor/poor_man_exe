package caddy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "http://localhost:2019"
	}
	return &Client{BaseURL: baseURL}
}

type Upstream struct {
	Dial string `json:"dial"`
}

type ReverseProxy struct {
	Handler   string     `json:"handler"`
	Upstreams []Upstream `json:"upstreams"`
}

type Match struct {
	Host []string `json:"host"`
}

type Route struct {
	Match  []Match        `json:"match"`
	Handle []ReverseProxy `json:"handle"`
}

func (c *Client) UpsertRoute(appName, domain string, port int) error {
	host := fmt.Sprintf("%s.%s", appName, domain)
	routeID := fmt.Sprintf("poor-exe-%s", appName)

	route := Route{
		Match: []Match{{Host: []string{host}}},
		Handle: []ReverseProxy{{
			Handler:   "reverse_proxy",
			Upstreams: []Upstream{{Dial: fmt.Sprintf("localhost:%d", port)}},
		}},
	}

	data, err := json.Marshal(route)
	if err != nil {
		return err
	}

	// Caddy's /config/ path allows IDs for easy upsert/delete
	url := fmt.Sprintf("%s/config/apps/http/servers/srv0/routes/%s", c.BaseURL, routeID)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("caddy api error: %s", resp.Status)
	}

	return nil
}

func (c *Client) DeleteRoute(appName string) error {
	routeID := fmt.Sprintf("poor-exe-%s", appName)
	url := fmt.Sprintf("%s/config/apps/http/servers/srv0/routes/%s", c.BaseURL, routeID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 404 is fine, means route doesn't exist
	if resp.StatusCode >= 400 && resp.StatusCode != 404 {
		return fmt.Errorf("caddy api error: %s", resp.Status)
	}

	return nil
}
