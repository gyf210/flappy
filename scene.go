package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
}

func NewScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background imags: %v", err)
	}

	b, err := NewBird(r)
	if err != nil {
		return nil, err
	}

	ps, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	return &scene{bg: bg, bird: b, pipes: ps}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	ch := make(chan error)

	go func() {
		defer close(ch)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				s.update()

				if s.bird.isDead() {
					drawTitle(r, "Game Over")
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					ch <- err
				}
			}
		}
	}()
	return ch
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.Event:
		return true
	case *sdl.MouseButtonEvent:
		s.bird.jump()
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent:
	default:
		log.Printf("unknown event %T", event)
	}
	return false
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}

	if err := s.pipes.paint(r); err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
}
