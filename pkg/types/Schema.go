package types

import "strings"

type Schema int

const (
	NONE Schema = iota
	HTTP
	HTTPS
)

func ParseSchema(schema string) Schema {
	schema = strings.ToLower(schema)
	switch schema {
	case "http":
		return HTTP
	case "https":
		return HTTPS
	default:
		return NONE
	}
}

func (p Schema) String() string {
	switch p {
	case HTTP:
		return "http"
	case HTTPS:
		return "https"
	default:
		return ""
	}
}

func (p Schema) ToSchemaString() string {
	switch p {
	case HTTP:
		return "http://"
	case HTTPS:
		return "https://"
	default:
		return ""
	}
}
