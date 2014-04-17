package main

import (
	"github.com/Centny/igtest"
	"os"
)

func main() {
	os.Exit(igtest.Run(os.Args[1:]))
}
