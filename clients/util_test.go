package clients

import (
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type UnmarshalStruct struct {
	Foo string `json:"foo"`
}

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		name    string
		resp    *http.Response
		want    UnmarshalStruct
		wantErr bool
	}{
		{
			name: "valid",
			resp: &http.Response{
				Body: io.NopCloser(strings.NewReader(`{"foo":"bar"}`)),
			},
			want:    UnmarshalStruct{Foo: "bar"},
			wantErr: false,
		},
		{
			name: "invalid json",
			resp: &http.Response{
				Body: io.NopCloser(strings.NewReader(`invalid`)),
			},
			want:    UnmarshalStruct{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := unmarshal[UnmarshalStruct](tc.resp)
			if (err != nil) != tc.wantErr {
				t.Errorf("unmarshal() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("unmarshal() = %v, want %v", got, tc.want)
			}
		})
	}
}

// // Mock response writer
// type mockResponseWriter struct {
// 	status  int
// 	written []byte
// }
//
// func (m *mockResponseWriter) Header() http.Header {
// 	return http.Header{}
// }
// func (m *mockResponseWriter) Write(p []byte) (int, error) {
// 	m.written = p
// 	return len(p), nil
// }
// func (m *mockResponseWriter) WriteHeader(status int) {
// 	m.status = status
// }
//
// func TestEncode(t *testing.T) {
//
// 	// Test data
// 	type response struct {
// 		Message string
// 	}
// 	data := response{Message: "Hello World"}
//
// 	wr := &mockResponseWriter{}
//
// 	// Request
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)
//
// 	// Call encode
// 	err := encode[response](wr, req, http.StatusOK, data)
//
// 	// Validate results
// 	if err != nil {
// 		t.Errorf("Unexpected error: %v", err)
// 	}
//
// 	if wr.status != http.StatusOK {
// 		t.Errorf("Unexpected status: got %v, want %v",
// 			wr.status, http.StatusOK)
// 	}
//
// 	expectContentType := "application/json"
// 	if wr.Header().Get("Content-Type") != expectContentType {
// 		t.Errorf("Unexpected Content-Type: got %v, want %v",
// 			wr.Header().Get("Content-Type"), expectContentType)
// 	}
//
// 	// Validate response body
// 	var resp response
// 	if err := json.Unmarshal(wr.written, &resp); err != nil {
// 		t.Errorf("Failed to unmarshal response: %v", err)
// 	}
//
// 	if resp.Message != data.Message {
// 		t.Errorf("Unexpected response body: got %v, want %v",
// 			resp.Message, data.Message)
// 	}
//
// }
