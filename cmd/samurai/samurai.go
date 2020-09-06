package main

import "github.com/swiftraptor/samurai/internal"

func main() {

	server := internal.MakeServer()
	server.Start("localhost", 9001)
}
