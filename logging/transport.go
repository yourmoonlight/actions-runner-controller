// Package logging provides various logging helpers for ARC
package logging

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/gregjones/httpcache"
)

const (
	// https://docs.github.com/en/rest/overview/resources-in-the-rest-api#rate-limiting
	headerRateLimitRemaining = "X-RateLimit-Remaining"
)

// Transport wraps a transport with metrics monitoring
type Transport struct {
	Transport http.RoundTripper

	Log *logr.Logger
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if resp != nil {
		t.log(req, resp)
	}
	return resp, err
}

func (t Transport) log(req *http.Request, resp *http.Response) {
	if t.Log == nil {
		return
	}

	var args []interface{}

	marked := resp.Header.Get(httpcache.XFromCache) == "1"

	args = append(args, "from_cache", marked, "method", req.Method, "url", req.URL.String())

	if !marked {
		// Do not log outdated rate limit remaining value

		remaining := resp.Header.Get(headerRateLimitRemaining)

		args = append(args, "ratelimit_remaining", remaining)
	}

	t.Log.V(3).Info("Seen HTTP response", args...)
}
