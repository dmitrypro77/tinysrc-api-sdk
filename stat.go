package tinysrc

import (
	"encoding/json"
	"net/http"
	"net/url"
	"tinysrc-api-sdk/models"
)

// Get List My Urls
func (client *Client) GetStatByHash(hash string, params models.StatRequest) (r *models.StatPaginatedResponse, errorResponse models.ErrorResponse) {
	values := url.Values{}

	values.Add("limit", string(rune(params.Limit)))
	values.Add("page", string(rune(params.Page)))
	values.Add("date-start", params.DateStart.Format(DATE_FORMAT))
	values.Add("date-end", params.DateEnd.Format(DATE_FORMAT))

	resp, e := client.sendRequest(http.MethodGet, "/client/stat/"+hash+"?"+values.Encode(), nil)
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	defer resp.Body.Close()

	if !client.isSuccess(resp.StatusCode) {
		apiErrors := client.parseErrorResponse(resp)

		return nil, *apiErrors
	}

	stat := models.StatPaginatedResponse{}
	e = json.NewDecoder(resp.Body).Decode(&stat)
	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	return &stat, errorResponse
}
