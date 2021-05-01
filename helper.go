package istiohelper

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	envoyReady = "http://localhost:15000/ready"
	istioReady = "http://localhost:15021/healthz/ready"

	//	envoyQuit = "http://localhost:15000/quitquitquit"
	istioQuit = "http://localhost:15020/quitquitquit"
)

// Helper object holds the state
type Helper struct {
	ok            bool
	debug         bool
	readyPort     string
	quitPort      string
	readyEndpoint string
	quitEndpoint  string
	readyAddr     string
	quitAddr      string
}

// ReadyPort to specify health check port
// 15000 - envoy
// 15020 - /healthz/ready
// 15021 - /quitquitquit
// If not set try everything
func ReadyPort(p string) func(*Helper) error {
	return func(d *Helper) error {
		d.readyPort = p
		return nil
	}
}

// ReadyEndpoint default "/healhtz/ready"
func ReadyEndpoint(r string) func(*Helper) error {
	return func(d *Helper) error {
		d.readyEndpoint = r
		return nil
	}
}

// QuitPort to specify the port to call quit
// 15000 - envoy
// 15021 - /quitquitquit
// Quitting Envoy may not be useful
func QuitPort(p string) func(*Helper) error {
	return func(d *Helper) error {
		d.quitPort = p
		return nil
	}
}

// QuitEndpoint default "/quitquitquit"
func QuitEndpoint(r string) func(*Helper) error {
	return func(d *Helper) error {
		d.readyEndpoint = r
		return nil
	}
}

// Debug will enable logging for debugging
var Debug func(d *Helper) error = func(d *Helper) error {
	d.debug = true
	return nil
}

// Wait for Istio (Envoy) proxy to report ready
func Wait(ok bool, options ...func(*Helper) error) *Helper {
	h := &Helper{
		ok: ok,
	}
	if !ok {
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
		h.readyAddr = fmt.Sprintf("http://localhost:%s%s", h.readyPort, h.readyEndpoint)
	}
	if h.quitAddr != "" {
		if h.quitEndpoint == "" {
			h.quitEndpoint = "/quitquitquit"
		}
		h.quitAddr = fmt.Sprintf("http://localhost:%s%s", h.quitPort, h.quitEndpoint)
	}

	for {
		if h.readyAddr != "" {
			if ok := h.checkReady(h.readyAddr); !ok {
				time.Sleep(time.Second)
				continue
			}
			return h
		}
		// No preference, try everything
		if ok := h.checkReady(istioReady); ok {
			return h
		}
		if ok := h.checkReady(envoyReady); ok {
			return h
		}
	}
}

func (h *Helper) checkReady(addr string) bool {
	resp, err := http.Get(addr)
	if err != nil {
		if h.debug {
			log.Printf("http.Get(%s) - %v\n", addr, err)
		}
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		if h.debug {
			log.Printf("http.Get(%s) - %v\n", addr, resp)
		}
		return false
	}
	return true
}

// Quit Istio (Envoy) proxy
func (h *Helper) Quit() {
	if !h.ok {
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
	if h.debug {
		log.Printf("http.Post(%s) - %v\n", h.quitAddr, resp)
	}
	defer resp.Body.Close()
}
