package istiohelper

import (
	"fmt"
	"net/http"
	"time"
)

const (
	envoyReady      = "http://localhost:15000/ready"
	istioOlderReady = "http://localhost:15020/healthz/ready"
	istioReady      = "http://localhost:15021/healthz/ready"

	//	envoyQuit = "http://localhost:15000/quitquitquit"
	istioQuit = "http://localhost:15020/quitquitquit"
)

// Helper object holds the state
type Helper struct {
	condition     bool
	legacy        bool
	logger        func(string)
	httpClient    *http.Client
	readyPort     string
	quitPort      string
	readyEndpoint string
	quitEndpoint  string
	readyAddr     string
	quitAddr      string
}

// Condition - pass a bool from a flag or config file
// to enable or disable istio checks
func Condition(t bool) func(*Helper) error {
	return func(h *Helper) error {
		h.condition = t
		return nil
	}
}

// HTTPClient can be used to set a custom http client for
// doing the Istio checks.
func HTTPClient(c *http.Client) func(*Helper) error {
	return func(h *Helper) error {
		if c == nil {
			return fmt.Errorf("http.Client cannot be nil")
		}
		h.httpClient = c
		return nil
	}
}

// Legacy - try other port and endpoint combinations that
// may help with older versions of Istio sidecars.
var Legacy = func(h *Helper) error {
	h.legacy = true
	return nil
}

// ReadyPort - specify a custom health check port.
// - 15000 - envoy
// - 15020 - /healthz/ready
// - 15021 - /quitquitquit
// If not set try everything
func ReadyPort(p string) func(*Helper) error {
	return func(h *Helper) error {
		h.readyPort = p
		return nil
	}
}

// ReadyEndpoint - specify a custom ready endpoint to check.
// default "/healhtz/ready"
func ReadyEndpoint(r string) func(*Helper) error {
	return func(h *Helper) error {
		h.readyEndpoint = r
		return nil
	}
}

// QuitPort - specify a custom port to call the quit endpoint.
// - 15000 Envoy
// - 15021 Istio
// Quitting Envoy may not be as useful as you might think
func QuitPort(p string) func(*Helper) error {
	return func(h *Helper) error {
		h.quitPort = p
		return nil
	}
}

// QuitEndpoint - specify a custom quit endpoint.
// default "/quitquitquit"
func QuitEndpoint(r string) func(*Helper) error {
	return func(h *Helper) error {
		h.readyEndpoint = r
		return nil
	}
}

// Logger - if a logger function is provided the health check
// will write basic trace information.
// To avoid dependencies on external packages,
// the log function is a basic `func(string)`. It is up to
// the user to display or write the string.
func Logger(logger func(string)) func(*Helper) error {
	return func(h *Helper) error {
		h.logger = logger
		return nil
	}
}

// Wait for Istio (Envoy) proxy to report ready
func Wait(options ...func(*Helper) error) *Helper {
	h := &Helper{
		condition: true,
		logger:    func(string) {},
	}
	if !h.condition {
		return h
	}

	for _, option := range options {
		if err := option(h); err != nil {
			return nil
		}
	}

	if h.readyPort != "" {
		if h.readyEndpoint == "" {
			h.readyEndpoint = "/ready"
		}
		h.readyAddr = "http://localhost:" + h.readyPort + ":" + h.readyEndpoint
	}
	if h.quitAddr != "" {
		if h.quitEndpoint == "" {
			h.quitEndpoint = "/quitquitquit"
		}
		h.quitAddr = "http://localhost:%s%s" + h.quitPort + h.quitEndpoint
	}

	for {
		if h.readyAddr != "" {
			if ok := h.checkReady(h.readyAddr); !ok {
				time.Sleep(time.Second)
				continue
			}
			return h
		}
		if ok := h.checkReady(istioReady); ok {
			return h
		}

		if h.legacy {
			if ok := h.checkReady(envoyReady); ok {
				return h
			}
			if ok := h.checkReady(istioOlderReady); ok {
				return h
			}
		}
		time.Sleep(time.Second)
	}
}

func (h *Helper) checkReady(addr string) bool {
	resp, err := http.Get(addr)
	if err != nil {
		h.logger("http.Get(" + addr + ") - " + err.Error())
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		h.logger("http.Get(" + addr + ") - " + resp.Status)
		return false
	}
	return true
}

// Quit Istio (Envoy) proxy
func (h *Helper) Quit() {
	if !h.condition {
		return
	}
	addr := istioQuit
	if h.quitAddr != "" {
		addr = h.quitAddr
	}
	resp, err := http.Post(addr, "application/json", nil)
	if err != nil {
		return
	}
	h.logger("http.Post(" + addr + ") - " + resp.Status)
	defer resp.Body.Close()
}
