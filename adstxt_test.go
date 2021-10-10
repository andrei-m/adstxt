package adstxt

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("CONTACT variable", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("CONTACT=foo"))
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Contact: []string{"foo"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("SUBDOMAIN variable", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("SUBDOMAIN=foo"))
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Subdomain: []string{"foo"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("INVENTORYPARTNERDOMAIN variable", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("INVENTORYPARTNERDOMAIN=foo"))
		assert.NoError(t, err)
		expected := map[Variable][]string{
			InventoryPartnerDomain: []string{"foo"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("case insensitive variable parsing", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("inventorypartnerdomain=foo"))
		assert.NoError(t, err)
		expected := map[Variable][]string{
			InventoryPartnerDomain: []string{"foo"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("skip unknown variables", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("foo=bar"))
		assert.NoError(t, err)
		assert.Len(t, adstxt.Variables, 0)
	})

	t.Run("value includes an additional '='", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("Contact=a=b"))
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Contact: []string{"a=b"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("multiple values are accumulated", func(t *testing.T) {
		rawAdsTxt := strings.NewReader(`contact=foo
contact=bar`)
		adstxt, err := Parse(rawAdsTxt)
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Contact: []string{"foo", "bar"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("full-line comment", func(t *testing.T) {
		rawAdsTxt := strings.NewReader(`#contact=foo
contact=bar`)
		adstxt, err := Parse(rawAdsTxt)
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Contact: []string{"bar"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("partial-line comment", func(t *testing.T) {
		rawAdsTxt := strings.NewReader(`contact=foo
subdomain=bar#comment`)
		adstxt, err := Parse(rawAdsTxt)
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Contact:   []string{"foo"},
			Subdomain: []string{"bar"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})

	t.Run("whitespace and empty lines are trimmed", func(t *testing.T) {
		rawAdsTxt := strings.NewReader(`
contact=foo

# another comment
subdomain=bar #comment`)
		adstxt, err := Parse(rawAdsTxt)
		assert.NoError(t, err)
		expected := map[Variable][]string{
			Contact:   []string{"foo"},
			Subdomain: []string{"bar"},
		}
		assert.Equal(t, expected, adstxt.Variables)
	})
}
