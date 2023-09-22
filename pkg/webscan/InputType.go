package webscan

type InputType uint8

const (
	Domain InputType = iota
	IPv4
	IPv6
)
