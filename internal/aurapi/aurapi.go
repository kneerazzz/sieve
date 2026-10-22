package aurapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURl = "https://aur.archlinux.org/rpc/v5"
	batchSize      = 100 //keeps the url under the lenth limit
)

type Info struct {
	Name           string  `json:"Name"`
	Version        string  `json:"Version"`
	Maintainer     *string `json:"Maintainer"`
	NumVotes       int     `json:"NumVotes"`
	Popularity     float64 `json:"Popularity"`
	OutOfDate      *int64  `json:"OutOfDate"`
	FirstSubmitted int64   `json:"FirstSubmitted"`
	LastModified   int64   `json:"LastModified"`
}

type Response struct {
	Type    string `json:"type"`
	Error   string `json:"error"`
	Results []Info `json:"results"`
}

type Client struct {
	BaseUrl string
	HTTP    *http.Client
}

func New() *Client {
	return &Client{
		BaseUrl: defaultBaseURl,
		HTTP:    &http.Client{Timeout: 15 * time.Second}, //Will give timeout after 15 second of trying
	}
}

// func will take client as argument Info is return type
// Info function looks up packages by name, Names Missing from the AUR are absent from the map
func (c *Client) Info(ctx context.Context, names []string) (map[string]Info, error) {
	out := make(map[string]Info, len(names))
	for start := 0; start < len(names); start += batchSize {
		end := min(start+batchSize, len(names))
		infos, err := c.fetch(ctx, names[start:end])
		if err != nil {
			return nil, err
		}
		for _, in := range infos {
			out[in.Name] = in
		}
	}
	return out, nil
}

//fetch function

func (c *Client) fetch(ctx context.Context, names []string) ([]Info, error) {
	q := url.Values{}
	for _, n := range names {
		q.Add("arg[]", n)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.BaseUrl+"/info?"+q.Encode(), nil)

	if err != nil {
		return nil, fmt.Errorf("Building request: %w", err)
	}
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling AUR RPC %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AUR RPC returned status %s", resp.Status)
	}
	var r Response

	if err := json.NewDecoder(io.LimitReader(resp.Body, 10<<20)).Decode(&r); err != nil {
		return nil, fmt.Errorf("Decoding AUR RPC response: %w", err)
	}
	if r.Type == "error" {
		return nil, fmt.Errorf("AUR RPC error: %s", r.Error)
	}
	return r.Results, nil
}
