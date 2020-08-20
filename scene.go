package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Scene struct {
	Triangles []*Triangle
	conf      *Conf
}

func (s *Scene) GenerateScene(conf *Conf) {
	s.conf = conf
	for i := 0; i < conf.Nb; i++ {
		s.CreateTriangle(getRandomPosition(s.conf.MaxX, s.conf.MaxY))
	}
}

func (s *Scene) CreateTriangle(position pixel.Vec) {

	min := 0.
	max := 1.

	triangle := &Triangle{
		A: Coordinate{
			Position: pixel.Vec{
				X: 0,
				Y: 0,
			},
		},
		B: Coordinate{
			Position: pixel.Vec{
				X: 10,
				Y: 0,
			},
		},
		C: Coordinate{
			Position: pixel.Vec{
				X: 5,
				Y: 10,
			},
		},
		G:                Coordinate{},
		Speed:            0,
		Direction:        pixel.Vec{X: 0, Y: 1},
		DirectionInitial: pixel.Vec{X: 0, Y: 1},
		Color:            pixel.RGB(min+rand.Float64()*(max-min), min+rand.Float64()*(max-min), min+rand.Float64()*(max-min)),
	}

	triangle.Speed = s.conf.MinSpeed + rand.Float64()*(s.conf.MaxSpeed-s.conf.MinSpeed)

	triangle.RefreshCenter()
	triangle.A.Position.X -= triangle.G.Position.X
	triangle.A.Position.Y -= triangle.G.Position.Y
	triangle.B.Position.X -= triangle.G.Position.X
	triangle.B.Position.Y -= triangle.G.Position.Y
	triangle.C.Position.X -= triangle.G.Position.X
	triangle.C.Position.Y -= triangle.G.Position.Y
	triangle.RefreshCenter()

	triangle.Translate(position)
	s.Triangles = append(s.Triangles, triangle)
	triangle.updateTriangleSprite()
}

func getRandomPosition(maxX, maxY float64) pixel.Vec {
	min := 0.

	randX := min + rand.Float64()*(maxX-min)
	randY := min + rand.Float64()*(maxY-min)
	return pixel.Vec{randX, randY}
}

func (s *Scene) dispatchPosition(vec pixel.Vec) {
	for _, triangle := range s.Triangles {
		triangle.UpdateDirection(vec)
	}
}

func (s *Scene) CatchEvent(win *pixelgl.Window) {
	s.conf.MaxX = win.Bounds().W()
	s.conf.MaxY = win.Bounds().H()
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		s.CreateTriangle(win.MousePosition())
		//s.dispatchPosition(win.MousePosition())
	}
}

func (s *Scene) Draw(win *pixelgl.Window) {
	for _, triangle := range s.Triangles {
		triangle.imd.Draw(win)
	}
}

func (s *Scene) Update() {
	for _, triangle := range s.Triangles {
		triangle.Move()
		triangle.Algo(s)
		triangle.Update()
		triangle.updateTriangleSprite()
	}
}
