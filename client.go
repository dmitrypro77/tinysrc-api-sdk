package tinysrc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dmitrypro77/tinysrc-api-sdk/models"
	"io"
	"net/http"
	"net/url"
)

// Client to send request to TinySRC API
type Client struct {
	httpClient *http.Client
	ApiKey     string
	baseURL    *url.URL
	ctx        context.Context
}

// Constructor of httpClient
func NewClient(ctx context.Context, apiKey string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	base, e := url.Parse(API_URL + fmt.Sprintf("/%s", VERSION))
	if e != nil {
		return nil, e
	}

	c := &Client{httpClient: httpClient, ApiKey: apiKey, baseURL: base, ctx: ctx}
	return c, nil
}

// Set required headers by TinySRC
func (client *Client) setRequestHeaders(req *http.Request) {
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-api-key", client.ApiKey)
}

// Send Request To TinySRC API
func (client *Client) sendRequest(method string, pathURL string, body io.Reader) (*http.Response, error) {
	rel, e := url.Parse(client.baseURL.Path + pathURL)
	if e != nil {
		return nil, e
	}

	fullURL := client.baseURL.ResolveReference(rel)

	req, e := http.NewRequest(method, fullURL.String(), body)
	if e != nil {
		return nil, e
	}

	client.setRequestHeaders(req)

	resp, e := client.do(req)
	if e != nil {
		return nil, e
	}

	return resp, nil
}

func (client *Client) do(req *http.Request) (*http.Response, error) {
	resp, e := client.httpClient.Do(req.WithContext(client.ctx))
	if e != nil {
		select {
		case <-client.ctx.Done():
			return nil, client.ctx.Err()
		default:
		}
		return nil, e
	}

	return resp, e
}

// Check If Response From Api Success
func (client *Client) isSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

// Check If Response From Api Success
func (client *Client) isUnauthorized(statusCode int) bool {
	return statusCode == 401
}

// Parse Response From Api
func (client *Client) parseErrorResponse(resp *http.Response) *models.ErrorResponse {
	errorResponse := models.ErrorResponse{}
	errorResponse.Status = resp.StatusCode

	if client.isUnauthorized(errorResponse.Status) {
		errorResponse.Errors = append(errorResponse.Errors, "Unauthorized")
	}

	e := json.NewDecoder(resp.Body).Decode(&errorResponse)

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return &errorResponse
	}

	return &errorResponse
}
