package tinysrc

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"
	"tinysrc-api-sdk/models"
)

func TestClient_do(t *testing.T) {
	type fields struct {
		Client *Client
	}

	type args struct {
		req *http.Request
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Test Passed")
	}))

	tsNotFound := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, _ = fmt.Fprint(w, "test 404")
	}))

	tsTimeout := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	defer ts.Close()
	defer tsNotFound.Close()
	defer tsTimeout.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/", nil)
	reqNotFound, _ := http.NewRequest(http.MethodGet, tsNotFound.URL+"/", nil)
	reqTimeout, _ := http.NewRequest(http.MethodGet, tsTimeout.URL+"/", nil)

	testClient, _ := NewClient(context.Background(), "test", nil)
	cont, cancel := context.WithTimeout(context.Background(), time.Duration(1)*time.Microsecond)
	testClientTimeout, _ := NewClient(cont, "test", nil)

	cancel()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			args: args{req: req},
			want: &http.Response{
				StatusCode: 200,
			},
			wantErr: false,
		},
		{
			name: "not_found_error",
			fields: fields{
				Client: testClient,
			},
			args: args{req: reqNotFound},
			want: &http.Response{
				StatusCode: 404,
			},
			wantErr: false,
		},
		{
			name: "timeout_error",
			fields: fields{
				Client: testClientTimeout,
			},
			args:    args{req: reqTimeout},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Client.do(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("do() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.StatusCode, tt.want.StatusCode) {
					t.Errorf("do() got = %v, want %v", got, tt.want)
				}
			} else {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("do() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestClient_isSuccess(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		statusCode int
	}

	testClient, _ := NewClient(context.Background(), "test", nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			args: args{statusCode: 200},
			want: true,
		},
		{
			name: "test_fail",
			fields: fields{
				Client: testClient,
			},
			args: args{statusCode: 200},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Client.isSuccess(tt.args.statusCode); got != tt.want {
				t.Errorf("isSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_isUnauthorized(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		statusCode int
	}

	testClient, _ := NewClient(context.Background(), "test", nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "test_success",
			fields: fields{
				Client: testClient,
			},
			args: args{statusCode: 401},
			want: true,
		},
		{
			name: "test_fail",
			fields: fields{
				Client: testClient,
			},
			args: args{statusCode: 200},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Client.isUnauthorized(tt.args.statusCode); got != tt.want {
				t.Errorf("isUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_parseErrorResponse(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		resp *http.Response
	}

	testClient, _ := NewClient(context.Background(), "test", nil)

	handler := func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "{\n\"validations\": {\n\"url\": [\n\" URL is not valid\"\n ]\n},\n\"errors\": [\n\"Validation Error Happened\"\n]\n}")
	}

	handlerNotError := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = io.WriteString(w, "{)")
	}

	req := httptest.NewRequest("GET", "http://example.com:8080", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	handlerNotError(w, req)

	resp := w.Result()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   *models.ErrorResponse
	}{
		{
			name: "test_validation_error",
			fields: fields{
				Client: testClient,
			},
			args: args{resp: resp},
			want: &models.ErrorResponse{
				Validations: map[string][]string{
					"url": {" URL is not valid"},
				},
				Errors: []string{"Validation Error Happened"},
				Status: 200,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Client.parseErrorResponse(tt.args.resp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseErrorResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_sendRequest(t *testing.T) {
	type fields struct {
		Client *Client
	}
	type args struct {
		method  string
		pathURL string
		body    io.Reader
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "Test Passed")
	}))

	defer ts.Close()

	testClient, _ := NewClient(context.Background(), "test", nil)

	testClient.baseURL = &url.URL{
		RawPath: ts.URL,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name:   "test_success",
			fields: fields{Client: testClient},
			args: args{
				method:  http.MethodGet,
				pathURL: ts.URL,
				body:    nil,
			},
			want:    200,
			wantErr: false,
		},
		{
			name:   "test_fail",
			fields: fields{Client: testClient},
			args: args{
				method:  http.MethodGet,
				pathURL: "*7/\\asdsad",
				body:    nil,
			},
			want:    500,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.fields.Client.sendRequest(tt.args.method, tt.args.pathURL, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !reflect.DeepEqual(got.StatusCode, tt.want) {
					t.Errorf("sendRequest() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestClient_setRequestHeaders(t *testing.T) {
	type fields struct {
		httpClient *http.Client
		ApiKey     string
		baseURL    *url.URL
		ctx        context.Context
	}

	req := httptest.NewRequest("GET", "http://example.com:8080", nil)

	type args struct {
		req *http.Request
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test_success",
			args: args{req: req},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				httpClient: tt.fields.httpClient,
				ApiKey:     tt.fields.ApiKey,
				baseURL:    tt.fields.baseURL,
				ctx:        tt.fields.ctx,
			}

			req, _ := http.NewRequest(http.MethodGet, "test", nil)
			client.setRequestHeaders(req)
			req.Header.Get("Accept")
			req.Header.Get("Content-Type")
			req.Header.Get("x-api-key")
		})
	}
}

func TestNewClient(t *testing.T) {
	type args struct {
		ctx        context.Context
		apiKey     string
		httpClient *http.Client
	}

	testClient, _ := NewClient(context.Background(), "test", nil)

	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "test_success",
			args: args{
				ctx:        context.Background(),
				apiKey:     "test",
				httpClient: http.DefaultClient,
			},
			want:    testClient,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.ctx, tt.args.apiKey, tt.args.httpClient)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}
