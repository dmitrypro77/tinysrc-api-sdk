package tinysrc

import (
	"context"
	"encoding/json"
	"github.com/dmitrypro77/tinysrc-api-sdk/models"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestClient_GetStatByHash(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		hash   string
		params models.StatRequest
	}

	successStat := models.StatPaginatedResponse{
		Data: []*models.StatResponse{
			{
				Ip:             "8.8.8.8",
				Bot:            false,
				Mobile:         false,
				Browser:        "ie",
				Os:             "mac",
				Platform:       "linux",
				Referer:        "test.com",
				BrowserVersion: "1.0",
				Created:        time.Time{},
			},
		},
		Total: 1,
	}

	successStatJson, _ := json.Marshal(&successStat)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, string(successStatJson))
	}))

	defer ts.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		Path: ts.URL,
	}
	successErr := models.ErrorResponse{}

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantR             *models.StatPaginatedResponse
		wantErrorResponse models.ErrorResponse
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			wantR:             &successStat,
			wantErrorResponse: successErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotErrorResponse := tt.fields.Client.GetStatByHash(tt.args.hash, tt.args.params)
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("GetStatByHash() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotErrorResponse, tt.wantErrorResponse) {
				t.Errorf("GetStatByHash() gotErrorResponse = %v, want %v", gotErrorResponse, tt.wantErrorResponse)
			}
		})
	}
}
