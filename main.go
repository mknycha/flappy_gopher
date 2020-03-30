package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return fmt.Errorf("could not initialize SDL: %w", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		return fmt.Errorf("could not initialize ttf: %w", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not initialize window or the rednerer: %w", err)
	}
	defer w.Destroy()
	_ = r

	err = drawTitle(r)
	if err != nil {
		return fmt.Errorf("could not draw a title: %w", err)
	}
	w.UpdateSurface()

	time.Sleep(1 * time.Second)

	if err := drawBackground(r); err != nil {
		return fmt.Errorf("could not draw background: %w", err)
	}

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}

	return nil
}

func drawBackground(r *sdl.Renderer) error {
	r.Clear()

	t, err := img.LoadTexture(r, "./res/images/background.png")
	if err != nil {
		return fmt.Errorf("could not load background: %w", err)
	}
	defer t.Destroy()
	err = r.Copy(t, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %w", err)
	}

	r.Present()
	return nil
}

func drawTitle(r *sdl.Renderer) error {
	r.Clear()

	f, err := ttf.OpenFont("./res/fonts/test.ttf", 20)
	if err != nil {
		return fmt.Errorf("error opening a font: %w", err)
	}
	defer f.Close()

	color := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := f.RenderUTF8Solid("Flappy Gopher", color)
	if err != nil {
		return fmt.Errorf("error rendering title: %w", err)
	}
	defer surface.Free()

	texture, err := r.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("error creating texture: %w", err)
	}
	defer texture.Destroy()

	err = r.Copy(texture, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy texture: %w", err)
	}
	r.Present()

	return nil
}
