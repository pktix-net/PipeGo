package proxy

import (
	"io"
	"log"
	"net"
	"slices"

	"pipego/internal/asn"
	"pipego/internal/route"
)

func TCPProxy(client net.Conn, r route.Route, asnDB *asn.DB) {
	defer client.Close()

	// ASN filter check
	if asnDB != nil {
		srcIP := extractIP(client.RemoteAddr())
		if srcIP != nil {
			asnNum, asnOrg := asnDB.Lookup(srcIP)

			// Whitelist mode
			if len(r.AllowASN) > 0 {
				if !slices.Contains(r.AllowASN, asnNum) {
					log.Printf("[INFO] asn blocked %s AS%d %s -> %s", client.RemoteAddr(), asnNum, asnOrg, r.Upstream)
					return
				}
			} else if len(r.DenyASN) > 0 {
				// Blacklist mode
				if slices.Contains(r.DenyASN, asnNum) {
					log.Printf("[INFO] asn denied %s AS%d %s -> %s", client.RemoteAddr(), asnNum, asnOrg, r.Upstream)
					return
				}
			}

			log.Printf("[INFO] asn pass %s AS%d %s -> %s", client.RemoteAddr(), asnNum, asnOrg, r.Upstream)
		}
	}

	server, err := net.Dial("tcp", r.Upstream)
	if err != nil {
		log.Printf("[ERROR] upstream connect failed: %s -> %s (%v)", client.RemoteAddr(), r.Upstream, err)
		return
	}
	defer server.Close()

	log.Printf(
		"[INFO] proxy %s <-> %s",
		client.RemoteAddr(),
		r.Upstream,
	)

	go io.Copy(server, client)
	io.Copy(client, server)
}

// extractIP extracts the IP from a net.Addr without the port
func extractIP(addr net.Addr) net.IP {
	switch a := addr.(type) {
	case *net.TCPAddr:
		return a.IP
	default:
		host, _, err := net.SplitHostPort(addr.String())
		if err != nil {
			return nil
		}
		return net.ParseIP(host)
	}
}
