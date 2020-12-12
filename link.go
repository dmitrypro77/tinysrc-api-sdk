package tinysrc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"tinysrc-api-sdk/models"
)

// Create a New Link
func (client *Client) CreateShortLink(requestData models.LinkRequest) (r *models.LinkResponse, errorResponse models.ErrorResponse) {
	body, e := json.Marshal(requestData)
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	resp, e := client.sendRequest("POST", "/create", bytes.NewBuffer(body))
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	defer resp.Body.Close()

	if !client.isSuccess(resp.StatusCode) {
		apiErrors := client.parseErrorResponse(resp)

		return nil, *apiErrors
	}

	linkResponse := models.LinkResponse{}
	e = json.NewDecoder(resp.Body).Decode(&linkResponse)
	if e != nil {
		return nil, errorResponse
	}

	return &linkResponse, errorResponse
}

// Get List My Urls
func (client *Client) GetListUrls(params models.ListUrlsRequest) (r *models.PaginatedLinkUserResponse, errorResponse models.ErrorResponse) {
	values := url.Values{}

	values.Add("limit", string(rune(params.Limit)))
	values.Add("page", string(rune(params.Page)))
	values.Add("query", params.Query)

	resp, e := client.sendRequest(http.MethodGet, "/client/url"+"?"+values.Encode(), nil)
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	defer resp.Body.Close()

	if !client.isSuccess(resp.StatusCode) {
		apiErrors := client.parseErrorResponse(resp)

		return nil, *apiErrors
	}

	listUrls := models.PaginatedLinkUserResponse{}
	e = json.NewDecoder(resp.Body).Decode(&listUrls)
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	return &listUrls, errorResponse
}

// Get List My Urls
func (client *Client) GetUrlByHash(hash string) (r *models.LinkUserResponse, errorResponse models.ErrorResponse) {
	resp, e := client.sendRequest(http.MethodGet, "/client/url/"+hash, nil)

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	defer resp.Body.Close()

	if !client.isSuccess(resp.StatusCode) {
		apiErrors := client.parseErrorResponse(resp)

		return nil, *apiErrors
	}

	urlInfo := models.LinkUserResponse{}
	e = json.NewDecoder(resp.Body).Decode(&urlInfo)
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	return &urlInfo, errorResponse
}

// Activate/Deactivate Hash
func (client *Client) SetActive(hash string, request *models.LinkActivationRequest) (status bool, errorResponse models.ErrorResponse) {
	body, e := json.Marshal(request)

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return false, errorResponse
	}

	resp, e := client.sendRequest(http.MethodPatch, "/client/"+hash, bytes.NewBuffer(body))

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return false, errorResponse
	}

	defer resp.Body.Close()

	if !client.isSuccess(resp.StatusCode) {
		apiErrors := client.parseErrorResponse(resp)

		return false, *apiErrors
	}

	urlInfo := models.LinkUserResponse{}

	e = json.NewDecoder(resp.Body).Decode(&urlInfo)

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return false, errorResponse
	}

	return true, errorResponse
}
