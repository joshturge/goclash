package goclash

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Client will send requests to the Clash of Clans API and serialize the response
type Client struct {
	BaseURL    *url.URL
	Token      string
	httpclient http.Client
	logger     *log.Logger

	Clan     *ClanService
	Player   *PlayerService
	League   *LeagueService
	Location *LocationService
	Label    *LabelService
}

type service struct {
	client *Client
}

// NewClient will create a new Client given a Clash of Clans API Token
func NewClient(token string) (*Client, error) {
	base, err := url.Parse("https://api.clashofclans.com/v1/")
	if err != nil {
		return nil, fmt.Errorf("could not pass base url: %s", err.Error())
	}

	client := Client{
		BaseURL: base,
		Token:   token,
		logger:  log.New(os.Stdout, "[LIBCLASH] ", 1),
	}

	commonSvc := service{&client}
	client.Clan = (*ClanService)(&commonSvc)
	client.Player = (*PlayerService)(&commonSvc)
	client.League = (*LeagueService)(&commonSvc)
	client.Location = (*LocationService)(&commonSvc)
	client.Label = (*LabelService)(&commonSvc)

	return &client, nil
}

// SetTimeout will set a timeout for requests
func (c *Client) SetTimeout(duration time.Duration) {
	c.httpclient.Timeout = duration
}

// NewRequest will create a new request to be sent to the Clash of Clans API
func (c *Client) NewRequest(path string, urlVal url.Values) (*http.Request, error) {
	var url strings.Builder
	url.WriteString(c.BaseURL.String())
	url.WriteString(path)
	if urlVal != nil || urlVal.Encode() != "" {
		url.WriteString("?")
		url.WriteString(urlVal.Encode())
	}

	c.logger.Println(url.String())

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create new request: %s", err.Error())
	}

	req.Header.Add("Accept", "application/json")

	var bearer strings.Builder
	bearer.WriteString("Bearer ")
	bearer.WriteString(c.Token)
	req.Header.Add("authorization", bearer.String())

	return req, nil
}

// Do will send a request to the Clash of Clans API and serialize the response into v
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpclient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not do request: %s", err.Error())
	}

	// TODO(joshturge): can't seem to get the content type of the response
	/*if resp.Header.Get("Content-Type") != "application/json" {
		return nil, fmt.Errorf("content type of response was not application/json")
	}*/

	var body bytes.Buffer
	if _, err = io.Copy(&body, resp.Body); err != nil {
		return nil, fmt.Errorf("could not copy response body: %s", err.Error())
	}
	resp.Body = ioutil.NopCloser(bytes.NewReader(body.Bytes()))

	decoder := json.NewDecoder(&body)

	if resp.StatusCode >= http.StatusBadRequest {
		errResp := ErrorResponse{Response: resp}
		if err = decoder.Decode(&errResp); err != nil {
			return nil, fmt.Errorf("could not decode response into an error: %s", err.Error())
		}

		return nil, fmt.Errorf("a server error has occurred: %s", errResp.Error())
	}

	if err = decoder.Decode(v); err != nil {
		return resp, fmt.Errorf("could not decode response body: %s", err.Error())
	}

	return resp, nil
}

// ErrorResponse is an error response from the Clash of Clans API
type ErrorResponse struct {
	Response *http.Response
	Reason   string `json:"reason"`
	Message  string `json:"message"`
	Type     string `json:"type"`
	Details  string `json:"details"`
}

// Error formats a string with information about an error
func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("[%s] %s: %d %s %+s",
		er.Response.Request.Method, er.Response.Request.URL.RequestURI(), er.Response.StatusCode,
		er.Message, er.Reason)
}
