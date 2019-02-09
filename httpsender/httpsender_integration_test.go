// +build integration

package httpsender_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/iulianclita/myhttp/httpsender"
)

func TestMake_integration(t *testing.T) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<h1>I love Go</h1>")
	}

	s := httptest.NewServer(http.HandlerFunc(fn))
	defer s.Close()

	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	tests := map[string]struct {
		url    string
		hash   string
		errNil bool
	}{
		"URL is reachable": {
			url:    s.URL,
			hash:   "5c0b86b58d53e45967b9105e84ace15f",
			errNil: true,
		},
		"URL is NOT reachable": {
			url:    fmt.Sprintf("%s1234567890", s.URL),
			hash:   "",
			errNil: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodGet, tc.url, nil)
			if err != nil {
				t.Fatalf("Failed creating request: %v", err)
			}

			want := tc.hash
			got, err := httpsender.Make(c, r)
			if got != want {
				t.Errorf("Make() = %s; want %s", got, want)
			}

			if tc.errNil {
				if err != nil {
					t.Errorf("Error is %v; want nil", err)
				}
			}

			if !tc.errNil {
				if err == nil {
					t.Errorf("Error is nil; want %v", err)
				}
			}
		})
	}

	r, err := http.NewRequest(http.MethodGet, s.URL, nil)
	if err != nil {
		t.Fatalf("Failed creating request: %v", err)
	}

	want := "5c0b86b58d53e45967b9105e84ace15f"
	got, err := httpsender.Make(c, r)

	if got != want {
		t.Errorf("Make() = %s; want %s", got, want)
	}

	if err != nil {
		t.Errorf("Error is %v; want nil", err)
	}
}
