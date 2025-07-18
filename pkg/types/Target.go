package types

import (
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"

	"log/slog"

	"github.com/miekg/dns"
)

type TargetType int

const (
	None TargetType = iota
	Domain
	Ipv4
	Ipv6
)

// isIPv4 checks if the given string is a valid IPv4 address
func isIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

// isIPv6 checks if the given string is a valid IPv6 address
func isIPv6(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() == nil
}

type Target struct {
	rawTarget  string
	schema     Schema
	parsedUrl  url.URL
	targetType TargetType
}

func NewTarget(targetString string) (Target, error) {
	var (
		err    error
		target = Target{
			rawTarget:  targetString,
			targetType: None,
			schema:     NONE,
		}
		parsedUrlPointer *url.URL
	)

	// Cleaning targetString
	targetString = strings.TrimSpace(targetString)
	targetString = strings.ToLower(targetString)

	target.schema = ParseSchema(strings.Split(targetString, "://")[0])

	switch target.schema {
	case HTTP, HTTPS:
		slog.Debug("target url contains schema", "schema", target.schema.String())
		parsedUrlPointer, err = url.Parse(targetString)
		if err != nil {
			break
		}
	case NONE:
		slog.Debug("No scheme provided in target, assuming 'https://' for the sake of parsing it")
		parsedUrlPointer, err = url.Parse(HTTPS.ToSchemaString() + targetString)
		if err != nil {
			break
		}
	default:
		return target, errors.New("unknown target schema")
	}

	if err != nil {
		slog.Error("Could not parse url", "error", err.Error())
		return target, err
	}

	target.parsedUrl = *parsedUrlPointer

	slog.Debug("target url parsed", "url", target.parsedUrl.String())

	switch {
	case net.ParseIP(target.parsedUrl.Hostname()) == nil: // If input is domain
		if _, ok := dns.IsDomainName(target.parsedUrl.Hostname()); ok {
			target.targetType = Domain
			slog.Debug("hostname identified as domain", "hostname", target.parsedUrl.Hostname())
		} else {
			target.targetType = None
			return target, errors.New("couldn't parse supposed ip address in hostname to neither ipv4 nor ipv6 address")
		}
	case isIPv4(target.parsedUrl.Hostname()):
		target.targetType = Ipv4
		slog.Debug("hostname identified as ipv4 address", "hostname", target.parsedUrl.Hostname())
	case isIPv6(target.parsedUrl.Hostname()):
		target.targetType = Ipv6
		slog.Debug("hostname identified as ipv6 address", "hostname", target.parsedUrl.Hostname())
	default:
		target.targetType = None
		return target, errors.New("couldn't parse supposed ip address in hostname to neither ipv4 nor ipv6 address")
	}

	return target, nil
}

func (target *Target) RawTarget() string {
	return target.rawTarget
}

func (target *Target) Schema() Schema {
	return target.schema
}

func (target *Target) OverrideSchema(schema Schema) {
	target.parsedUrl.Scheme = schema.String()
	target.schema = schema
}

func (target *Target) TargetType() TargetType {
	return target.targetType
}

func (target *Target) ParsedUrl() url.URL {
	return target.parsedUrl
}

func (target *Target) UrlString() string {
	return target.parsedUrl.String()
}

// This should not be used, to ensure it's not edited by mistake
// func (target *Target) ParsedUrlPointer() *url.URL {
// 	return &target.parsedUrl
// }

// Hostname without port
func (target *Target) Hostname() string {
	return target.parsedUrl.Hostname()
}

// Overrides the hostname of the target (no port)
func (target *Target) OverrideHostname(hostname string) {
	target.parsedUrl.Host = hostname + ":" + target.parsedUrl.Port()
}

// Hostname with port
// <hostname>:<port> or <hostname> if no port is specified
func (target *Target) Host() string {
	return target.parsedUrl.Host
}

// Port without hostname
func (target *Target) Port() string {
	return target.parsedUrl.Port()
}

// Parses the port as uint16, returns 0 if parsing fails
func (target *Target) PortAsUint16() uint16 {
	port, err := strconv.ParseUint(target.parsedUrl.Port(), 10, 16)
	if err != nil {
		return 0
	}
	return uint16(port)
}

func (target *Target) OverridePort(port string) {
	target.parsedUrl.Host = target.parsedUrl.Hostname() + ":" + port
}

func (target *Target) Path() string {
	return target.parsedUrl.EscapedPath()
}

// TODO add tests for each case, including error cases (invalid domain, ...)
