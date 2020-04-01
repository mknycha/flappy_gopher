package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type bird struct {
	// We can have many goroutines accessing mutex, but if anyone is writing, none can access
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture

	x, y  int32
	w, h  int32
	speed float64
	dead  bool
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
	return &bird{time: 0, textures: textures, y: 300, x: 10, w: 50, h: 43, speed: 0}, nil
}

func (b *bird) update() {
	// read only lock
	b.mu.RLock()
	defer b.mu.RUnlock()
	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {

	// Little hack that changes the bird frame only ten times per second
	i := b.time / 10 % len(b.textures)
	// This way our bird will be in the middle of the screen
	rect := &sdl.Rect{X: b.x, Y: (600 - b.y) - b.h/2, W: b.w, H: b.h}
	err := r.Copy(b.textures[i], nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy background: %w", err)
	}
	return nil
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = -jumpSpeed
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.x > b.x+b.w { // too far right
		return
	}
	if p.x+p.w < b.x { // too far left
		return
	}
	if !p.inverted && p.h < b.y-b.h/2 { // pipe too low
		return
	}
	if p.inverted && (600-p.h) > b.y-b.h/2 { // inverted pipe too high
		return
	}
	b.dead = true
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, t := range b.textures {
		t.Destroy()
	}
}
