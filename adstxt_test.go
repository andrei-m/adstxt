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

	t.Run("DIRECT relationship record", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("foo,bar,DIRECT"))
		assert.NoError(t, err)
		expected := []Record{
			{
				AdSystemDomain:  "foo",
				SellerAccountID: "bar",
				Relationship:    Direct,
			},
		}
		assert.Equal(t, expected, adstxt.Records)
	})

	t.Run("url-encoded comma in SellerAccountID", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("foo,foo%2Cbar,DIRECT"))
		assert.NoError(t, err)
		expected := []Record{
			{
				AdSystemDomain:  "foo",
				SellerAccountID: "foo,bar",
				Relationship:    Direct,
			},
		}
		assert.Equal(t, expected, adstxt.Records)
	})

	t.Run("RESELLER relationship record", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("foo,bar,RESELLER"))
		assert.NoError(t, err)
		expected := []Record{
			{
				AdSystemDomain:  "foo",
				SellerAccountID: "bar",
				Relationship:    Reseller,
			},
		}
		assert.Equal(t, expected, adstxt.Records)
	})

	t.Run("unknown relationship", func(t *testing.T) {
		_, err := Parse(strings.NewReader("foo,bar,baz"))
		assert.Equal(t, errUnrecognizedRelationshipType, err)
	})

	t.Run("missing ad system domain", func(t *testing.T) {
		_, err := Parse(strings.NewReader("foo,,DIRECT"))
		assert.Equal(t, errNoSellerAccountID, err)
	})

	t.Run("missing seller account ID", func(t *testing.T) {
		_, err := Parse(strings.NewReader(",foo,DIRECT"))
		assert.Equal(t, errNoAdSystemDomain, err)
	})

	t.Run("record with certification authority ID", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("foo,bar,RESELLER,baz"))
		assert.NoError(t, err)
		expected := []Record{
			{
				AdSystemDomain:  "foo",
				SellerAccountID: "bar",
				Relationship:    Reseller,
				CertAuthorityID: "baz",
			},
		}
		assert.Equal(t, expected, adstxt.Records)
	})

	t.Run("file with all features", func(t *testing.T) {
		rawAdsTxt := strings.NewReader(`
# comment
foo,bar,DIRECT,baz
one,two,RESELLER

# another comment
contact=foo
contact=foobar
subdomain=bar #comment`)
		adstxt, err := Parse(rawAdsTxt)
		assert.NoError(t, err)
		expected := AdsTxt{
			Records: []Record{
				{
					AdSystemDomain:  "foo",
					SellerAccountID: "bar",
					Relationship:    Direct,
					CertAuthorityID: "baz",
				},
				{
					AdSystemDomain:  "one",
					SellerAccountID: "two",
					Relationship:    Reseller,
				},
			},
			Variables: map[Variable][]string{
				Contact:   []string{"foo", "foobar"},
				Subdomain: []string{"bar"},
			},
		}
		assert.Equal(t, expected, adstxt)
	})

	t.Run("placeholder", func(t *testing.T) {
		adstxt, err := Parse(strings.NewReader("placeholder.example.com, placeholder, DIRECT, placeholder"))
		assert.NoError(t, err)
		assert.Empty(t, adstxt.Records)
	})
}

func Test_Record_isPlaceholder(t *testing.T) {
	record := Record{}
	assert.False(t, record.isPlaceholder())
	placeholder := Record{
		AdSystemDomain:  "placeholder.example.com",
		SellerAccountID: "placeholder",
	}
	assert.True(t, placeholder.isPlaceholder())
}
