package adstxt

type AdsTxt struct {
	Records   []Record
	Variables map[Variable][]string
}

type Record struct {
	AdSystemDomain  string
	SellerAccountID string
	Relationship    RelationshipType
	CertAuthorityID string
}

type RelationshipType int

const (
	RelationshipTypeUnspecified RelationshipType = iota
	Direct
	Reseller
)

type Variable int

const (
	VariableUnspecified Variable = iota
	Contact
	Subdomain
	InventoryPartnerDomain
)
