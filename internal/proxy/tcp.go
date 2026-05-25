package proxy

import (
	"io"
	"log"
	"net"
)

func TCPProxy(client net.Conn, upstream string) {
	defer client.Close()

	server, err := net.Dial("tcp", upstream)
	if err != nil {
		log.Printf("[ERROR] upstream connect failed: %s -> %s (%v)", client.RemoteAddr(), upstream, err)
		return
	}
	defer server.Close()

	log.Printf(
		"[INFO] proxy %s <-> %s",
		client.RemoteAddr(),
		upstream,
	)

	go io.Copy(server, client)
	io.Copy(client, server)
}
