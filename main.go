package main

import (
	"github.com/knightazura/infrastructure"
)

func main() {
	// Bootstrap the application
	infrastructure.Bootstrap()

	// Dispatch the app
	infrastructure.Dispatch()
}