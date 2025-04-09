package client

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/ubombar/soa/api"
)

type ZoteroClientConfig struct {
	Enpoint *url.URL
}

const DefaultZoteroClientEndpoint = "http://localhost:23119/better-bibtex/"

// This client uses the Zotero's Bette rBibtext plugin, ensure it is installed
// to reach the endpoint
type ZoteroClient struct {
	client *http.Client
	cfg    *ZoteroClientConfig
}

func NewZoteroClient(cfg *ZoteroClientConfig) (*ZoteroClient, error) {
	client := &http.Client{}

	if cfg == nil {
		u, err := url.Parse(DefaultZoteroClientEndpoint)
		if err != nil {
			return nil, err
		}
		cfg = &ZoteroClientConfig{
			Enpoint: u,
		}
	}

	return &ZoteroClient{
		cfg:    cfg,
		client: client,
	}, nil
}

// This invokes the selection UI of Zotero bibtext plugin.
func (c *ZoteroClient) SelectBibTextEntries() ([]api.ZoteroCitationEntry, error) {
	enpointURL := c.cfg.Enpoint.ResolveReference(&url.URL{Path: "cayw"})
	q := enpointURL.Query()
	q.Set("format", "json")
	enpointURL.RawQuery = q.Encode()

	req, err := http.NewRequest("GET", enpointURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var citationEntries []api.ZoteroCitationEntry
	if err := json.NewDecoder(resp.Body).Decode(&citationEntries); err != nil {
		return nil, err
	}

	return citationEntries, nil
}

// Gets the attachement from a citationKey
func (c *ZoteroClient) GetAttachements(citationKey string) ([]api.ZoteroAttachementItem, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "item.attachments",
		"params":  []string{citationKey},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	enpointURL := c.cfg.Enpoint.ResolveReference(&url.URL{Path: "json-rpc"})

	req, err := http.NewRequest("POST", enpointURL.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var attachementResponse api.ZoteroAttachementResponse
	if err := json.NewDecoder(resp.Body).Decode(&attachementResponse); err != nil {
		return nil, err
	}

	return attachementResponse.Result, nil
}
