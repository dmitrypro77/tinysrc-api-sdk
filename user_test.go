package tinysrc

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"tinysrc-api-sdk/models"
)

func TestClient_GetCurrentUser(t *testing.T) {
	type fields struct {
		Client *Client
	}

	successUser := models.CurrentUserResponse{
		Username: "test",
		ApiKey:   "test",
		Active:   1,
		Plan:     1,
		Banned:   0,
		Email:    "test@test.com",
	}

	successUserJson, _ := json.Marshal(&successUser)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(successUserJson))
	}))

	defer ts.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		Path: ts.URL,
	}

	testClientFail, _ := NewClient(context.Background(), "test", nil)

	testClientFail.baseURL = &url.URL{
		Path: "//---%2f",
	}

	successErr := models.ErrorResponse{}

	tests := []struct {
		name              string
		fields            fields
		wantR             *models.CurrentUserResponse
		wantErrorResponse models.ErrorResponse
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			wantR:             &successUser,
			wantErrorResponse: successErr,
		},
		{
			name: "test_fail",
			fields: fields{
				Client: testClientFail,
			},
			wantR: nil,
			wantErrorResponse: models.ErrorResponse{
				Errors: []string{"parse \"//---%2f/client/user\": invalid URL escape \"%2f\""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotErrorResponse := tt.fields.Client.GetCurrentUser()
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetCurrentUser() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotErrorResponse, tt.wantErrorResponse) {
				t.Errorf("GetCurrentUser() gotErrorResponse = %v, want %v", gotErrorResponse, tt.wantErrorResponse)
			}
		})
	}
}
