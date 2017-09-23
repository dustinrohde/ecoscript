package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	mapfile, err := ParseMapfile("examples/Mapfile")
	if err != nil {
		log.Fatal(err)
	}

	world := mapfile.ToWorld()

	for {
		fmt.Println(world.Layer(0).Display())
		world.Tick()
		time.Sleep(500 * time.Millisecond)
	}
}
