package certificate

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

//AboutToExpire returns true if time.Now() + earlyExpirity is after cert.NotAfter date.
func AboutToExpire(certPEM []byte, earlyExpirity string) (expiring bool, date time.Time, err error) {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return false, time.Now(), fmt.Errorf("failed to decode pem")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false, time.Now(), fmt.Errorf("failed to parse certificate. %s", err)
	}
	duration, err := time.ParseDuration(earlyExpirity)
	if err != nil {
		return false, time.Now(), fmt.Errorf("failed to parse duration from string %s. %s", earlyExpirity, err)
	}
	if time.Now().Add(duration).After(cert.NotAfter) {
		return true, cert.NotAfter, nil
	}
	return false, cert.NotAfter, nil
}
