package tinysrc

import (
	"encoding/json"
	"github.com/dmitrypro77/tinysrc-api-sdk/models"
	"net/http"
)

// Get Current User Information
func (client *Client) GetCurrentUser() (r *models.CurrentUserResponse, errorResponse models.ErrorResponse) {
	resp, e := client.sendRequest(http.MethodGet, "/client/user", nil)

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	defer resp.Body.Close()

	if !client.isSuccess(resp.StatusCode) {
		apiErrors := client.parseErrorResponse(resp)

		return nil, *apiErrors
	}

	currentUserResponse := models.CurrentUserResponse{}
	e = json.NewDecoder(resp.Body).Decode(&currentUserResponse)

	if e != nil {
		errorResponse.Errors = append(errorResponse.Errors, e.Error())
		return nil, errorResponse
	}

	return &currentUserResponse, errorResponse
}
