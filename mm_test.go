package mm_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k1LoW/httpstub"
	"github.com/k1LoW/mm"
)

type testBuilder struct {
	buildFunc func(req *http.Request) (func(next http.Handler) http.Handler, bool)
}

func (b *testBuilder) Middleware(req *http.Request) (func(next http.Handler) http.Handler, bool) {
	return b.buildFunc(req)
}

func TestMM(t *testing.T) {
	// testMw is a test middleware
	// Set "X-Test" header to all requests
	testMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Test", "test")
			next.ServeHTTP(w, r)
		})
	}

	tests := []struct {
		name     string
		builders []mm.Builder
		req      *http.Request
		want     *http.Response
	}{
		{
			name:     "No middleware",
			builders: nil,
			req:      &http.Request{Method: http.MethodGet, URL: mustParseURL("http://example.com")},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": {"text/plain; charset=utf-8"}},
			},
		},
		{
			name: "Set testMw to all requests",
			builders: []mm.Builder{
				&testBuilder{
					buildFunc: func(req *http.Request) (func(next http.Handler) http.Handler, bool) {
						return testMw, true
					},
				},
			},
			req: &http.Request{Method: http.MethodGet, URL: mustParseURL("http://example.com")},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"text/plain; charset=utf-8"},
					"X-Test":       {"test"},
				},
			},
		},
		{
			name: "Set testMw to only GET requests (1/2)",
			builders: []mm.Builder{
				&testBuilder{
					buildFunc: func(req *http.Request) (func(next http.Handler) http.Handler, bool) {
						if req.Method != http.MethodGet {
							return nil, false
						}
						return testMw, true
					},
				},
			},
			req: &http.Request{Method: http.MethodHead, URL: mustParseURL("http://example.com")},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"text/plain; charset=utf-8"},
				},
			},
		},
		{
			name: "Set testMw to only GET requests (2/2)",
			builders: []mm.Builder{
				&testBuilder{
					buildFunc: func(req *http.Request) (func(next http.Handler) http.Handler, bool) {
						if req.Method != http.MethodGet {
							return nil, false
						}
						return testMw, true
					},
				},
			},
			req: &http.Request{Method: http.MethodGet, URL: mustParseURL("http://example.com")},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": {"text/plain; charset=utf-8"},
					"X-Test":       {"test"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httpstub.NewRouter(t)
			r.Match(func(r *http.Request) bool { return true }).Response(http.StatusOK, []byte("Hello"))
			m := mm.New(tt.builders...)
			ts := httptest.NewServer(m(r))
			tu := mustParseURL(ts.URL)
			t.Cleanup(ts.Close)
			tc := ts.Client()
			tt.req.URL.Scheme = tu.Scheme
			tt.req.URL.Host = tu.Host
			got, err := tc.Do(tt.req)
			if err != nil {
				t.Fatal(err)
			}
			opts := []cmp.Option{
				cmpopts.IgnoreFields(http.Response{}, "Status", "Proto", "ProtoMajor", "ProtoMinor", "ContentLength", "TransferEncoding", "Uncompressed", "Trailer", "Request", "Close", "Body"),
			}
			// header ignore fields
			got.Header.Del("Content-Length")
			got.Header.Del("Date")
			if diff := cmp.Diff(tt.want, got, opts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func mustParseURL(urlstr string) *url.URL {
	u, err := url.Parse(urlstr)
	if err != nil {
		panic(err) //nostyle:dontpanic
	}
	return u
}
