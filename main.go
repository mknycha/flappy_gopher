package main

import (
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()
}
