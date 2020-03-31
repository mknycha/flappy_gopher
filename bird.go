package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type bird struct {
	time     int
	textures []*sdl.Texture
	y, speed float64
}

const (
	gravity   = 0.25
	jumpSpeed = 5
)

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i <= 4; i++ {
		path := fmt.Sprintf("./res/images/frame-%v.png", i)
		t, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load background: %w", err)
		}
		textures = append(textures, t)
	}
	return &bird{time: 0, textures: textures, y: 300, speed: 0}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	if b.y < 0 {
		b.speed = -b.speed
		b.y = 0
	}
	b.speed += gravity

	// Little hack that changes the bird frame only ten times per second
	i := b.time / 10 % len(b.textures)
	// This way our bird will be in the middle of the screen
	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}
	err := r.Copy(b.textures[i], nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy background: %w", err)
	}
	return nil
}

func (b *bird) jump() {
	b.speed = -jumpSpeed
}

func (b *bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}
