package webscan

type InputType uint8

const (
	Domain InputType = iota
	DomainWithPath
	IPv4
	IPv6
)
