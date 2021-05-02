package istiohelper_test

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/physcat/istiohelper"
)

func ExampleCondition() {
	wait := flag.Bool("wait-for-istio", false, "wait for Istio to start")
	flag.Parse()

	defer istiohelper.Wait(istiohelper.Condition(*wait)).Quit()
}

func ExampleHTTPClient() {
	c := http.Client{
		Timeout: time.Second,
	}
	defer istiohelper.Wait(istiohelper.HTTPClient(&c)).Quit()
}

func ExampleWait() {
	defer istiohelper.Wait().Quit()
}

func ExampleLegacy() {
	defer istiohelper.Wait(istiohelper.Legacy).Quit()
}

func ExampleReadyPort() {
	defer istiohelper.Wait(istiohelper.ReadyPort("15000")).Quit()
}

func ExampleLogger() {
	defer istiohelper.Wait(istiohelper.Logger(func(msg string) { log.Println(msg) })).Quit()
}
