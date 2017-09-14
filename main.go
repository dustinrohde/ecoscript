package main

import (
	"github.com/davecgh/go-spew/spew"
)

func main() {
	mapfile := ParseMapfile("examples/Mapfile")
	spew.Dump(mapfile)
	err := mapfile.Sanitize()
	spew.Dump(err)
}

