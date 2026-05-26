package listener

import (
	"log"
	"net"

	"pipego/internal/asn"
	"pipego/internal/proxy"
	"pipego/internal/route"
)

func Start(r route.Route, asnDB *asn.DB) {
	ln, err := net.Listen(r.Protocol, r.Listen)
	if err != nil {
		log.Printf("[ERROR] listen failed: %s (%v)", r.Listen, err)
		return
	}

	log.Printf(
		"[INFO] listening %s %s -> %s",
		r.Protocol,
		r.Listen,
		r.Upstream,
	)

	Serve(ln, r, asnDB)
}

func Serve(ln net.Listener, r route.Route, asnDB *asn.DB) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}

		go proxy.TCPProxy(conn, r, asnDB)
	}
}
