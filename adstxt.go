package adstxt

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"
	"unicode"
)

type AdsTxt struct {
	Records   []Record
	Variables map[Variable][]string
}

func (a AdsTxt) String() string {
	recordSuffix := ""
	if len(a.Records) != 1 {
		recordSuffix = "s"
	}
	variableSuffix := ""
	if len(a.Variables) != 1 {
		variableSuffix = "s"
	}
	return fmt.Sprintf("%d record%s, %d variable%s", len(a.Records), recordSuffix, len(a.Variables), variableSuffix)
}

type Record struct {
	AdSystemDomain  string
	SellerAccountID string
	Relationship    RelationshipType
	CertAuthorityID string
	Extension       string
}

func (r Record) isPlaceholder() bool {
	return r.AdSystemDomain == "placeholder.example.com" && r.SellerAccountID == "placeholder"
}

type RelationshipType int

const (
	RelationshipTypeUnspecified RelationshipType = iota
	Direct
	Reseller
)

var strToRelationshipType = map[string]RelationshipType{
	"DIRECT":   Direct,
	"RESELLER": Reseller,
}

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

var (
	errNoRecord                     = errors.New("no record")
	errNoAdSystemDomain             = errors.New("ad system domain is required")
	errNoSellerAccountID            = errors.New("seller account ID is required")
	errUnrecognizedRelationshipType = errors.New("unrecognized relationship type is neither DIRECT or RESELLER")
)

func parseRecord(line string) (Record, error) {
	extension := ""
	extensionSeparatorIdx := strings.Index(line, ";")
	if extensionSeparatorIdx != -1 {
		extension = line[extensionSeparatorIdx+1:]
		line = line[:extensionSeparatorIdx]
	}
	recordSplit := strings.Split(line, ",")
	if len(recordSplit) < 3 {
		return Record{}, errNoRecord
	}
	relationship := strToRelationshipType[strings.ToUpper(recordSplit[2])]
	if relationship == RelationshipTypeUnspecified {
		return Record{}, errUnrecognizedRelationshipType
	}

	decodedAdSystemDomain, err := url.QueryUnescape(recordSplit[0])
	if err != nil {
		return Record{}, err
	}
	if len(decodedAdSystemDomain) == 0 {
		return Record{}, errNoAdSystemDomain
	}
	decodedSellerAccountId, err := url.QueryUnescape(recordSplit[1])
	if err != nil {
		return Record{}, err
	}
	if len(decodedSellerAccountId) == 0 {
		return Record{}, errNoSellerAccountID
	}

	record := Record{
		AdSystemDomain:  decodedAdSystemDomain,
		SellerAccountID: decodedSellerAccountId,
		Relationship:    relationship,
		Extension:       extension,
	}
	if len(recordSplit) > 3 {
		decodedCertAuthorityID, err := url.QueryUnescape(recordSplit[3])
		if err != nil {
			return Record{}, err
		}
		record.CertAuthorityID = decodedCertAuthorityID
	}
	return record, nil
}

func processComment(line string) string {
	idx := strings.Index(line, "#")
	if idx == -1 {
		return line
	}
	return line[:idx]
}

func stripWhitespace(line string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, line)
}

// Parse parses ads.txt from in
func Parse(in io.Reader) (AdsTxt, error) {
	variables := map[Variable][]string{}
	records := []Record{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		line = processComment(line)
		line = stripWhitespace(line)
		if len(line) == 0 {
			continue
		}

		record, err := parseRecord(line)
		if err == errNoRecord {
			// no-op; continue to parse a variable
		} else if err != nil {
			return AdsTxt{}, err
		} else if record.isPlaceholder() {
			continue
		} else {
			records = append(records, record)
		}

		variable, value := parseVariable(line)
		if variable != VariableUnspecified {
			variables[variable] = append(variables[variable], value)
		}
	}
	if err := scanner.Err(); err != nil {
		return AdsTxt{}, err
	}
	return AdsTxt{Records: records, Variables: variables}, nil
}
