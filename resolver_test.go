package adstxt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extractRootDomain(t *testing.T) {
	assert.Equal(t, "foo.com", extractRootDomain("bar.foo.com"))
	assert.Equal(t, "foo.co.uk", extractRootDomain("bar.foo.co.uk"))
	assert.Equal(t, "foo.com", extractRootDomain("subdomain.foo.com"))
}
