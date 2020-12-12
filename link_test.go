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

func TestClient_CreateShortLink(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		requestData models.LinkRequest
	}

	successLink := models.LinkResponse{
		Url:          "http://test.com",
		StatUrl:      "http://test.com/stat/122",
		StatPassword: "123",
		Password:     "123",
		AuthRequired: 1,
	}

	successLinkJson, _ := json.Marshal(&successLink)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(successLinkJson))
	}))

	defer ts.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		Path: ts.URL,
	}

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantR             *models.LinkResponse
		wantErrorResponse models.ErrorResponse
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			wantR:             &successLink,
			wantErrorResponse: models.ErrorResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotErrorResponse := tt.fields.Client.CreateShortLink(tt.args.requestData)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("CreateShortLink() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotErrorResponse, tt.wantErrorResponse) {
				t.Errorf("CreateShortLink() gotErrorResponse = %v, want %v", gotErrorResponse, tt.wantErrorResponse)
			}
		})
	}
}

func TestClient_GetListUrls(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		params models.ListUrlsRequest
	}

	success := models.PaginatedLinkUserResponse{
		Data: []*models.LinkUserResponse{
			{
				Url:      "test",
				Hash:     "test",
				Active:   1,
				Clicks:   1,
				Bots:     1,
				Password: "test",
			},
		},
		Total: 1,
	}

	successJson, _ := json.Marshal(&success)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(successJson))
	}))

	defer ts.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		Path: ts.URL,
	}

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantR             *models.PaginatedLinkUserResponse
		wantErrorResponse models.ErrorResponse
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			args: args{params: models.ListUrlsRequest{
				Limit: 10,
				Page:  1,
				Query: "test",
			}},
			wantR:             &success,
			wantErrorResponse: models.ErrorResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotErrorResponse := tt.fields.Client.GetListUrls(tt.args.params)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetListUrls() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotErrorResponse, tt.wantErrorResponse) {
				t.Errorf("GetListUrls() gotErrorResponse = %v, want %v", gotErrorResponse, tt.wantErrorResponse)
			}
		})
	}
}

func TestClient_GetUrlByHash(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		hash string
	}

	success := models.LinkUserResponse{
		Url:          "test.com",
		Hash:         "test",
		AuthRequired: 0,
		Password:     "test",
		StatPassword: "test",
		QRCode:       "test",
		Active:       1,
		Clicks:       11,
		Bots:         11,
		StatUrl:      "http://test.com",
	}

	successJson, _ := json.Marshal(&success)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(successJson))
	}))

	defer ts.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		Path: ts.URL,
	}

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantR             *models.LinkUserResponse
		wantErrorResponse models.ErrorResponse
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			wantR:             &success,
			wantErrorResponse: models.ErrorResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotErrorResponse := tt.fields.Client.GetUrlByHash(tt.args.hash)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetUrlByHash() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotErrorResponse, tt.wantErrorResponse) {
				t.Errorf("GetUrlByHash() gotErrorResponse = %v, want %v", gotErrorResponse, tt.wantErrorResponse)
			}
		})
	}
}

func TestClient_SetActive(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		hash    string
		request *models.LinkActivationRequest
	}

	urlInfo := models.LinkUserResponse{
		Url:    "test",
		Active: 1,
		Clicks: 11,
		Bots:   12,
	}

	fail := models.ErrorResponse{
		Validations: map[string][]string{"url": {"url is not valid"}},
		Status:      404,
	}

	successJson, _ := json.Marshal(&urlInfo)
	failJson, _ := json.Marshal(&fail)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(successJson))
	}))

	ts404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = io.WriteString(w, string(failJson))
	}))

	defer ts.Close()

	defer ts404.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		Path: ts.URL,
	}

	testClient404, _ := NewClient(context.Background(), "test", nil)

	testClient404.baseURL = &url.URL{
		Path: ts404.URL,
	}

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantStatus        bool
		wantErrorResponse models.ErrorResponse
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			args: args{
				hash:    "test",
				request: &models.LinkActivationRequest{Active: true},
			},
			wantStatus:        true,
			wantErrorResponse: models.ErrorResponse{},
		},
		{
			name: "test_fail",
			fields: fields{
				Client: testClient404,
			},
			args: args{
				hash:    "test",
				request: &models.LinkActivationRequest{Active: true},
			},
			wantStatus:        false,
			wantErrorResponse: fail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotErrorResponse := tt.fields.Client.SetActive(tt.args.hash, tt.args.request)
			if gotStatus != tt.wantStatus {
				t.Errorf("SetActive() gotStatus = %v, want %v", gotStatus, tt.wantStatus)
			}
			if !reflect.DeepEqual(gotErrorResponse, tt.wantErrorResponse) {
				t.Errorf("SetActive() gotErrorResponse = %v, want %v", gotErrorResponse, tt.wantErrorResponse)
			}
		})
	}
}
