package adstxt

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractRootDomain(t *testing.T) {
	assert.Equal(t, "foo.com", extractRootDomain("bar.foo.com"))
	assert.Equal(t, "foo.co.uk", extractRootDomain("bar.foo.co.uk"))
	assert.Equal(t, "foo.com", extractRootDomain("subdomain.foo.com"))
}

func Test_CheckRedirect(t *testing.T) {
	t.Run("stop following after 10 redirects", func(t *testing.T) {
		via := make([]*http.Request, 10)
		for i := range via {
			via[i] = httptest.NewRequest("GET", "https://www.example.com", nil)
		}
		assert.Equal(t, http.ErrUseLastResponse, CheckRedirect(via[0], via))
	})

	t.Run("redirecting more than once within the same root domain is allowed", func(t *testing.T) {
		via := []*http.Request{
			httptest.NewRequest("GET", "https://foo.example.com", nil),
			httptest.NewRequest("GET", "https://bar.example.com", nil),
		}
		assert.NoError(t, CheckRedirect(httptest.NewRequest("GET", "https://baz.example.com", nil), via))
	})

	t.Run("one external redirect", func(t *testing.T) {
		via := []*http.Request{
			httptest.NewRequest("GET", "https://foo.example.com", nil),
		}
		assert.NoError(t, CheckRedirect(httptest.NewRequest("GET", "https://baz.thirdparty.net", nil), via))
	})

	t.Run("multiple external redirects", func(t *testing.T) {
		via := []*http.Request{
			httptest.NewRequest("GET", "https://foo.example.com", nil),
			httptest.NewRequest("GET", "https://baz.thirdparty.net", nil),
		}
		assert.Equal(t, errMultipleExternalRedirects, CheckRedirect(httptest.NewRequest("GET", "https://foobar.thirdparty.net", nil), via))
	})

	t.Run("internal redirect followed by one external redirect", func(t *testing.T) {
		via := []*http.Request{
			httptest.NewRequest("GET", "https://example.com", nil),
			httptest.NewRequest("GET", "https://www.example.com", nil),
		}
		assert.NoError(t, CheckRedirect(httptest.NewRequest("GET", "https://baz.thirdparty.com", nil), via))
	})

	t.Run("internal redirect followed by two external redirects", func(t *testing.T) {
		via := []*http.Request{
			httptest.NewRequest("GET", "https://example.com", nil),
			httptest.NewRequest("GET", "https://www.example.com", nil),
			httptest.NewRequest("GET", "https://foobar.thirdparty.net", nil),
		}
		assert.Equal(t, errMultipleExternalRedirects, CheckRedirect(httptest.NewRequest("GET", "https://baz.thirdparty.net", nil), via))
	})
}
