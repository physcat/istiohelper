package istiohelper_test

import (
	"fmt"

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
		istiohelper.Port("15000"),
		istiohelper.Debug).Quit()
	// Output: Not waiting for Istio proxy
}
