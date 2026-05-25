package listener

import (
	"log"
	"net"

	"pipego/internal/proxy"
	"pipego/internal/route"
)

func Start(r route.Route) {
	listener, err := net.Listen(r.Protocol, r.Listen)
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

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[ERROR] accept failed: %v", err)
			continue
		}

		go proxy.TCPProxy(conn, r.Upstream)
	}
}
