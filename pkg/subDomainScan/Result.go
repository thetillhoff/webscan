package subDomainScan

type Result struct {
	subdomainsFromTlsScan map[string]struct{}
	subdomainsFromCrtSh   map[string]struct{}
}
