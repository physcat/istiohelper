package istiohelper

import (
	"fmt"
	"net/http"
	"time"
)

// Helper object holds the state
type Helper struct {
	ok            bool
	port          string
	readyEndpoint string
	quitEndpoint  string
	readyAddr     string
	quitAddr      string
}

// Port option may be needed to set the port to 15000 for older
// versions of Istio
func Port(p string) func(*Helper) error {
	return func(d *Helper) error {
		d.port = p
		return nil
	}
}

// Wait for Istio (Envoy) proxy to report ready
func Wait(ok bool, options ...func(*Helper) error) *Helper {
	h := &Helper{
		ok:            ok,
		port:          "15020",
		readyEndpoint: "/ready",
		quitEndpoint:  "/quitquitquit",
	}
	if !ok {
		return h
	}

	for _, option := range options {
		if err := option(h); err != nil {
			return nil
		}
	}

	h.readyAddr = fmt.Sprintf("http://localhost:%s%s", h.port, h.readyEndpoint)
	h.quitAddr = fmt.Sprintf("http://localhost:%s%s", h.port, h.quitEndpoint)

	for {
		resp, err := http.Get(h.readyAddr)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			time.Sleep(time.Second)
			continue
		}
		return h
	}
}

// Quit Istio (Envoy) proxy
func (h *Helper) Quit() {
	if !h.ok {
		return
	}
	resp, err := http.Post(h.quitAddr, "application/json", nil)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}
