package istiohelper_test

import (
	"fmt"
	"log"

	"github.com/physcat/istiohelper"
)

func ExampleWait() {
	fmt.Println("Not waiting for Istio proxy")
	defer istiohelper.Wait(false).Quit()
	// Output: Not waiting for Istio proxy
}

func ExampleWait_withPort() {
	fmt.Println("Not waiting for Istio proxy")
	defer istiohelper.Wait(false,
		istiohelper.ReadyPort("15000"),
		istiohelper.Logger(func(msg string) { log.Println(msg) }),
	).Quit()
	// Output: Not waiting for Istio proxy
}
