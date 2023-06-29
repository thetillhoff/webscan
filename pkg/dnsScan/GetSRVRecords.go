package dnsScan

// TODO

// import (
// 	"context"
// 	"net"
// )

// func (engine Engine) GetSRVRecords(url string) (Engine, error) {
// 	var (
// 		err error

// 		records    []string
// 		srvRecords []*net.NS
// 	)

// 	srvRecords, err = engine.resolver.LookupSRV(context.Background(), url)
// 	if err, ok := err.(*net.DNSError); ok && err.IsNotFound {
// 		// No SRV record available
// 	} else if err != nil {
// 		return engine, err
// 	}

// 	for _, record := range srvRecords {
// 		records = append(records, record.Host)
// 	}

// 	engine.srvRecords = records
// 	return engine, nil
// }
