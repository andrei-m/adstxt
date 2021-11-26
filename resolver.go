package adstxt

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var (
	errRequestFailed             = errors.New("ads.txt request was not successful")
	errMultipleExternalRedirects = errors.New("at most one redirect to a destination outside the original root domain is allowed")
)

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
	if !isHTTPSuccess(resp) {
		return AdsTxt{}, err
	}
	defer resp.Body.Close()
	return Parse(resp.Body)
}

func isHTTPSuccess(resp *http.Response) bool {
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// CheckRedirect is a http.Client CheckRedirect function that implements at-most-one external redirect as specified in section 3.1 of the ads.txt spec (v. 1.03)
func CheckRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return http.ErrUseLastResponse
	}
	sameRootDomain := extractRootDomain(via[0].URL.Host) == extractRootDomain(req.URL.Host)
	if !sameRootDomain && len(via) > 1 {
		previousRedirectSameRootDomain := extractRootDomain(via[0].URL.Host) == extractRootDomain(via[len(via)-1].URL.Host)
		if !previousRedirectSameRootDomain {
			return errMultipleExternalRedirects
		}
	}
	return nil
}

func extractRootDomain(host string) string {
	suffix, _ := publicsuffix.PublicSuffix(host)
	suffixElements := strings.Split(suffix, ".")
	hostElements := strings.Split(host, ".")
	lastPreSuffixElementIdx := len(hostElements) - len(suffixElements) - 1
	if lastPreSuffixElementIdx >= 0 {
		return strings.Join(append([]string{hostElements[lastPreSuffixElementIdx]}, suffixElements...), ".")
	}
	return suffix
}
