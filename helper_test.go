package istiohelper_test

import (
	"log"
	"net/http"
	"time"

	"github.com/physcat/istiohelper"
)

func ExampleHTTPClient() {
	c := http.Client{
		Timeout: time.Second,
	}
	defer istiohelper.Wait(true, istiohelper.HTTPClient(&c)).Quit()
}

func ExampleWait() {
	defer istiohelper.Wait(true).Quit()
}

func ExampleLegacy() {
	defer istiohelper.Wait(true, istiohelper.Legacy).Quit()
}

func ExampleReadyPort() {
	defer istiohelper.Wait(true, istiohelper.ReadyPort("15000")).Quit()
}

func ExampleLogger() {
	defer istiohelper.Wait(false,
		istiohelper.Logger(func(msg string) { log.Println(msg) }),
	).Quit()
}
