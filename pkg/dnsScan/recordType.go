package dnsScan

type recordType int

const (
	NS recordType = iota
	A
	AAAA
	CNAME
	TXT
)
