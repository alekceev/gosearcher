package search

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateMidPrices(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	tests := []struct {
		name    string
		search  string
		urls    []string
		wantW   []string
		wantErr bool
	}{
		{
			name:    "1",
			search:  "USD",
			urls:    []string{srv.URL + "/rest/v2/alpha/us"},
			wantW:   []string{srv.URL + "/rest/v2/alpha/us"},
			wantErr: false,
		},
		{
			name:    "2",
			search:  "USD",
			urls:    []string{srv.URL + "/rest/v2/alpha/ru", srv.URL + "/rest/v2/alpha/us"},
			wantW:   []string{srv.URL + "/rest/v2/alpha/us"},
			wantErr: false,
		},
		{
			name:    "3",
			search:  "USD",
			urls:    []string{srv.URL + "/rest/v2/alpha/ru"},
			wantW:   []string{},
			wantErr: false,
		},
		{
			name:    "4",
			search:  "USD",
			urls:    []string{"err"},
			wantW:   []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotW, err := SearchInUrls(tt.search, tt.urls)
			if (err != nil) != tt.wantErr {
				fmt.Println(err)
				t.Errorf("SearchInUrls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotW) != len(tt.wantW) {
				t.Errorf("SearchInUrls() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/rest/v2/alpha/us", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"currency":"USD"}`))
	})
	handler.HandleFunc("/rest/v2/alpha/ru", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"currency":"RUB"}`))
	})

	srv := httptest.NewServer(handler)

	return srv
}
