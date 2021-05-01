# Istio helper function

When using an Istio proxy it may be useful to wait for the
proxy to be ready before trying to connect to other
services like databases.

```go
package main

import (
	"fmt"

	"github.com/physcat/istiohelper"
)

func main() {
    fmt.Println("Waiting for Istio sidecar")
    defer istiohelper.Wait(true).Quit()
    fmt.Println("Istio sidecar is ready")


    // the defer will signal the sidecar
    // to quit at the end of main
}

```
