package adstxt

import (
	"bufio"
	"io"
	"strings"
)

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

var strToVariable = map[string]Variable{
	"CONTACT":                Contact,
	"SUBDOMAIN":              Subdomain,
	"INVENTORYPARTNERDOMAIN": InventoryPartnerDomain,
}

func parseVariable(line string) (Variable, string) {
	keyValSplit := strings.Split(line, "=")
	if len(keyValSplit) > 1 {
		key := keyValSplit[0]
		for str, variable := range strToVariable {
			if strings.EqualFold(str, key) {
				// account for possible additional '=' symbols in the arbitrary value
				return variable, strings.Join(keyValSplit[1:], "=")
			}
		}
	}
	return VariableUnspecified, ""
}

func Parse(in io.Reader) (AdsTxt, error) {
	variables := map[Variable][]string{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		variable, value := parseVariable(line)
		if variable != VariableUnspecified {
			variables[variable] = append(variables[variable], value)
		}
	}
	if err := scanner.Err(); err != nil {
		return AdsTxt{}, err
	}
	return AdsTxt{Variables: variables}, nil
}
