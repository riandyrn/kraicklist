package main

import (
	"github.com/knightazura/infrastructure"
)

func main() {
	// Load configuration
	infrastructure.Bootstrap()

	// Dispatch the app
	infrastructure.Dispatch()
}