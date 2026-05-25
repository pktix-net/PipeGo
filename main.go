package main

import (
	"log"
	"pipego/internal/config"
	"pipego/internal/listener"
)

func main() {
	log.Println("[INFO] PipeGo starting...")

	routes, err := config.LoadRoutes("pipego.routes")
	if err != nil {
		log.Fatalf("[FATAL] load routes failed: %v", err)
	}

	if len(routes) == 0 {
		log.Fatal("[FATAL] no routes found")
	}

	for _, r := range routes {
		go listener.Start(r)
	}

	select {}
}
