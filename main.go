package main

import (
	"flag"
	"log"

	"pipego/internal/manager"
	"pipego/internal/web"
)

func main() {
	webPort := flag.Int("p", 12138, "Web management port")
	webPass := flag.String("a", "123456", "Web management password")
	asnDBPath := flag.String("asn-db", "GeoLite2-ASN.mmdb", "MaxMind GeoLite2-ASN database path (set empty to disable)")
	flag.Parse()

	log.Println("[INFO] PipeGo starting...")

	mgr := manager.New(*asnDBPath)

	// Try to load initial routes; warn but don't exit if file is missing (can be created via web)
	if err := mgr.LoadAndStart("pipego.routes"); err != nil {
		log.Printf("[WARN] load routes failed: %v (you can manage routes via web interface)", err)
	}

	if len(mgr.Routes()) == 0 {
		log.Println("[WARN] no routes loaded, waiting for web management configuration...")
	}

	// Start web management server
	webServer := web.New(mgr, "admin", *webPass, *webPort)
	go webServer.Start()

	select {}
}
