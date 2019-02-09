package httpsender_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iulianclita/myhttp/httpsender"
)

var ErrNotReachable = errors.New("URL not reachable")

type fakeClient struct {
	reachableURL bool
}

func (fc fakeClient) Do(req *http.Request) (*http.Response, error) {
	if fc.reachableURL {
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("I love Go")),
		}, nil
	}

	return nil, ErrNotReachable
}

func TestMake(t *testing.T) {
	tests := map[string]struct {
		reachableURL bool
		hash         string
		err          error
	}{
		"URL is reachable": {
			reachableURL: true,
			hash:         "3beabce759f9cefc252844a87d3601b2",
			err:          nil,
		},
		"URL is NOT reachable": {
			reachableURL: false,
			hash:         "",
			err:          ErrNotReachable,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fc := fakeClient{
				reachableURL: tc.reachableURL,
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			hash, err := httpsender.Make(fc, req)
			if hash != tc.hash {
				t.Errorf("Make() = %s; want %s", hash, tc.hash)
			}
			if err != tc.err {
				t.Errorf("Error is %v; want %v", err, tc.err)
			}
		})
	}
}
