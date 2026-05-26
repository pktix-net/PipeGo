package manager

import (
	"log"
	"net"
	"sync"

	"pipego/internal/asn"
	"pipego/internal/config"
	"pipego/internal/listener"
	"pipego/internal/route"
)

type Manager struct {
	mu        sync.Mutex
	routes    []route.Route
	listeners []net.Listener
	asnDB     *asn.DB
}

// New creates a Manager. If asnDBPath is non-empty, loads the ASN database.
func New(asnDBPath string) *Manager {
	m := &Manager{}
	if asnDBPath != "" {
		db, err := asn.Open(asnDBPath)
		if err != nil {
			log.Printf("[WARN] asn db open failed: %v, asn filtering disabled", err)
		} else {
			m.asnDB = db
			log.Printf("[INFO] asn db loaded: %s", asnDBPath)
		}
	}
	return m
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

		rt := r // capture loop variable
		log.Printf(
			"[INFO] listening %s %s -> %s",
			rt.Protocol,
			rt.Listen,
			rt.Upstream,
		)

		go listener.Serve(ln, rt, m.asnDB)
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
