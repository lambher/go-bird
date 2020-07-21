package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Scene struct {
	Triangles []*Triangle
	maxX      float64
	maxY      float64
}

func (s *Scene) GenerateScene(maxX, maxY float64) {
	s.maxX = maxX
	s.maxY = maxY
	nb := 300
	for i := 0; i < nb; i++ {
		triangle := s.CreateTriangle()

		triangle.updateTriangleSprite()
		s.Triangles = append(s.Triangles, triangle)
	}
}

func (s *Scene) CreateTriangle() *Triangle {

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

	min = 0.5
	max = 3.

	triangle.Speed = min + rand.Float64()*(max-min)

	triangle.RefreshCenter()
	triangle.A.Position.X -= triangle.G.Position.X
	triangle.A.Position.Y -= triangle.G.Position.Y
	triangle.B.Position.X -= triangle.G.Position.X
	triangle.B.Position.Y -= triangle.G.Position.Y
	triangle.C.Position.X -= triangle.G.Position.X
	triangle.C.Position.Y -= triangle.G.Position.Y
	triangle.RefreshCenter()

	min = 0.
	max = s.maxX

	randX := min + rand.Float64()*(max-min)
	max = s.maxY
	randY := min + rand.Float64()*(max-min)

	triangle.Translate(pixel.Vec{randX, randY})

	return triangle
}

func (s *Scene) dispatchPosition(vec pixel.Vec) {
	for _, triangle := range s.Triangles {
		triangle.UpdateDirection(vec)
	}
}

func (s *Scene) CatchEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		s.dispatchPosition(win.MousePosition())
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
