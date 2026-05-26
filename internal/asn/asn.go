package asn

import (
	"fmt"
	"net"

	"github.com/oschwald/geoip2-golang"
)

// DB wraps a MaxMind GeoLite2-ASN database reader
type DB struct {
	reader *geoip2.Reader
}

// Open opens a MaxMind ASN .mmdb file
func Open(path string) (*DB, error) {
	reader, err := geoip2.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open asn db: %w", err)
	}
	return &DB{reader: reader}, nil
}

// Lookup returns the ASN number and organization name for the given IP
// Returns 0 and "" if no ASN data found
func (db *DB) Lookup(ip net.IP) (int, string) {
	record, err := db.reader.ASN(ip)
	if err != nil {
		return 0, ""
	}
	return int(record.AutonomousSystemNumber), record.AutonomousSystemOrganization
}

// Close releases the database resources
func (db *DB) Close() {
	db.reader.Close()
}
