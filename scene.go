package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg   *sdl.Texture
	bird *bird
	pipe *pipe
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

	pipe, err := newPipe(r)
	if err != nil {
		return nil, fmt.Errorf("could initialize pipe: %w", err)
	}

	return &scene{bg: bg, bird: bird, pipe: pipe}, nil
}

// returns a channel we want to read from
func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case event := <-events:
				if done := s.handleEvent(event); done {
					return
				}
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game over")
					time.Sleep(1 * time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch e := event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.KeyboardEvent:
		log.Println(e.Type)
		if e.Type == uint32(768) {
			s.bird.jump()
		}
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent:
		// just to clean logs
	default:
		log.Printf("unknown event: %T", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipe.update()
	s.bird.touch(s.pipe)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipe.restart()
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %w", err)
	}

	err = s.bird.paint(r)
	if err != nil {
		return fmt.Errorf("could not paint bird: %w", err)
	}
	err = s.pipe.paint(r)
	if err != nil {
		return fmt.Errorf("could not paint pipe: %w", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipe.destroy()
}
