package adstxt

import (
	"errors"
	"net/http"
	"net/url"
)

var errRequestFailed = errors.New("ads.txt request was not successful")

// Resolve requests and parses ads.txt from the provided host
func Resolve(doer Doer, host string) (AdsTxt, error) {
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/ads.txt",
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return AdsTxt{}, err
	}
	resp, err := doer.Do(req)
	if err != nil {
		return AdsTxt{}, errRequestFailed
	}
	if resp.StatusCode == 404 {
		return AdsTxt{}, nil
	}
	//TODO: handle at most one redirect
	if !isHTTPSuccess(resp) {
		return AdsTxt{}, err
	}
	defer resp.Body.Close()
	return Parse(resp.Body)
}

func isHTTPSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}
