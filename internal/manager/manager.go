package manager

import (
	"log"
	"net"
	"sync"

	"pipego/internal/config"
	"pipego/internal/listener"
	"pipego/internal/route"
)

type Manager struct {
	mu        sync.Mutex
	routes    []route.Route
	listeners []net.Listener
}

func New() *Manager {
	return &Manager{}
}

func (m *Manager) Routes() []route.Route {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]route.Route, len(m.routes))
	copy(out, m.routes)
	return out
}

func (m *Manager) LoadAndStart(path string) error {
	routes, err := config.LoadRoutes(path)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopAllLocked()

	for _, r := range routes {
		ln, err := net.Listen(r.Protocol, r.Listen)
		if err != nil {
			log.Printf("[ERROR] listen failed: %s (%v)", r.Listen, err)
			continue
		}

		m.listeners = append(m.listeners, ln)

		log.Printf(
			"[INFO] listening %s %s -> %s",
			r.Protocol,
			r.Listen,
			r.Upstream,
		)

		go listener.Serve(ln, r.Upstream)
	}

	m.routes = routes
	return nil
}

func (m *Manager) stopAllLocked() {
	for _, ln := range m.listeners {
		ln.Close()
	}
	m.listeners = nil
	m.routes = nil
}
