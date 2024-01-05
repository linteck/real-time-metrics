//go:build js && wasm

package main

import (
	"webjs/gojsdomframe"
)

func main() {
	// webpages.CreateFrame(&wc)
	// webpagesjs.CreateFrame(&wc)
	gojsdomframe.CreateFrame()
	// If you don’t do that select{}, or something to that effect,
	// you’ll find this error in the Developer console of your browser:
	// Error: Go program has already exited
	select {}
}
