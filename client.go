package microcms

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	serviceDomain  string
	apiKey         string
	globalDraftKey string
}

type ClientParams func(*Client)

type Params struct {
	contentId string
}

type RequestParams func(*Params)

func CreateClient(serviceDomain, apiKey string, params ...ClientParams) *Client {
	c := &Client{
		serviceDomain:  serviceDomain,
		apiKey:         apiKey,
		globalDraftKey: "",
	}
	for _, param := range params {
		param(c)
	}
	return c
}

func (c *Client) makeRequest(method, endpoint, contentId string) (*http.Request, error) {
	url := createUrl(c.serviceDomain, endpoint, contentId)

	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-KEY", c.apiKey)

	if c.globalDraftKey != "" {
		req.Header.Set("X-GLOBAL-DRAFT-KEY", c.globalDraftKey)
	}

	return req, nil
}

func (c *Client) Get(endpoint string, params ...RequestParams) (interface{}, error) {
	p := &Params{
		contentId: "",
	}

	for _, params := range params {
		params(p)
	}

	req, err := c.makeRequest(GET, endpoint, p.contentId)
	res, _ := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var data interface{}
	if err := parseBody(res, &data); err != nil {
		return nil, err
	}

	return data, err
}

func parseBody(res *http.Response, v interface{}) error {
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(v)
}

func createUrl(serviceDomain, endpoint, contentId string) string {
	base := fmt.Sprintf("https://%s.%s/api/%s/%s", serviceDomain, BASE_DOMAIN, API_VERSION, endpoint)
	if contentId != "" {
		base := fmt.Sprintf("%s/%s", base, contentId)
		return base
	}

	return base
}

func GlobalDraftKey(v string) ClientParams {
	return func(c *Client) {
		c.globalDraftKey = v
	}
}

func ContentId(v string) RequestParams {
	return func(p *Params) {
		p.contentId = v
	}
}
