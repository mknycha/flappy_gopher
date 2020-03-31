package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time int
	bg   *sdl.Texture
	bird *bird
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "./res/images/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background: %w", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, fmt.Errorf("could initialize bird: %w", err)
	}

	return &scene{bg: bg, bird: bird}, nil
}

// returns a channel we want to read from
func (s *scene) run(events chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case event := <-events:
				switch event.(type) {
				case *sdl.QuitEvent:
					return
				}
			// we try to get context done
			// context will be done whenever the scene should finish
			case <-tick:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++

	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %w", err)
	}

	err = s.bird.paint(r)
	if err != nil {
		return fmt.Errorf("could not paint bird: %w", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
}
