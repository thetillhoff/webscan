package webscan

import (
	"errors"
	"net"
	"net/url"
	"strings"

	"golang.org/x/exp/slog"
)

type Target struct {
	rawTarget string
	parsedUrl *url.URL
	isDomain  bool
	isIpv4    bool
	isIpv6    bool
}

func NewTarget(targetString string) (Target, error) {
	var (
		err    error
		target = Target{
			rawTarget: targetString,
			isDomain:  false,
			isIpv4:    false,
			isIpv6:    false,
		}
	)

	target.parsedUrl, err = url.Parse(targetString)
	if err != nil {
		slog.Error("Could not parse url", "error", err.Error())
		return target, err
	}

	if target.parsedUrl.Hostname() == "" {
		slog.Debug("No scheme provided in target url, assuming 'https://'")
		target.parsedUrl, err = url.ParseRequestURI("https://" + targetString)
		if err != nil {
			slog.Error("Could not parse url, even with assumed schema", "error", err.Error())
			return target, err
		}
	}
	slog.Debug("target url parsed")

	if net.ParseIP(target.parsedUrl.Hostname()) == nil { // If input is domain
		target.isDomain = true
		slog.Debug("hostname is domain")
	} else {
		slog.Debug("hostname is ip address")
		if strings.Count(target.parsedUrl.Hostname(), ".") == 3 && len(target.parsedUrl.Hostname()) >= 7 {
			target.isIpv4 = true
			slog.Debug("hostname is ipv4 address")
		} else if strings.Count(target.parsedUrl.Hostname(), ":") >= 2 && len(target.parsedUrl.Hostname()) >= 3 {
			target.isIpv6 = true
			slog.Debug("hostname is ipv6 address")
		} else {
			return target, errors.New("couldn't parse supposed ip address in hostname to neither ipv4 nor ipv6 address")
		}
	}

	return target, nil
}

func (target *Target) IsDomain() bool {
	return target.isDomain
}

func (target *Target) IsIpv4() bool {
	return target.isIpv4
}

func (target *Target) IsIpv6() bool {
	return target.isIpv6
}

func (target *Target) Hostname() string {
	return target.parsedUrl.Hostname()
}

func (target *Target) Port() string {
	return target.parsedUrl.Port()
}

func (target *Target) Path() string {
	return target.parsedUrl.EscapedPath()
}

// TODO add tests for each case, including error cases (invalid domain, ...)
