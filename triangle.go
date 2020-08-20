package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
)

type Coordinate struct {
	Position    pixel.Vec
	Translation pixel.Vec
}

func (t *Triangle) Rotate(teta float64) {
	t.A.Position = t.A.Position.Rotated(teta)
	t.B.Position = t.B.Position.Rotated(teta)
	t.C.Position = t.C.Position.Rotated(teta)
	t.G.Position = t.G.Position.Rotated(teta)
}

func (t *Triangle) GetPos() pixel.Vec {
	return t.G.Position.Add(t.G.Translation)
}

func (t *Triangle) GetDistance(triangle *Triangle) float64 {
	return triangle.GetPos().Sub(t.GetPos()).Len()
}

func (t *Triangle) GetColorDiff(triangle *Triangle) float64 {
	diff := math.Abs(t.Color.R - triangle.Color.R)
	diff += math.Abs(t.Color.G - triangle.Color.G)
	diff += math.Abs(t.Color.B - triangle.Color.B)
	return diff / 3
}

type Triangle struct {
	imd   *imdraw.IMDraw
	Color pixel.RGBA

	A Coordinate
	B Coordinate
	C Coordinate

	G Coordinate

	Speed            float64
	Direction        pixel.Vec
	DirectionInitial pixel.Vec

	Followers  []*Triangle
	Subscriber *Triangle
}

func (t *Triangle) Move() {
	min := -.2
	max := .2

	r := min + rand.Float64()*(max-min)
	t.Direction = t.Direction.Rotated(r)
}

func (t *Triangle) Update() {
	t.ApplyTranslation()
	t.RefreshCenter()
	t.ApplyRotation()
}

func (t *Triangle) Algo(scene *Scene) {
	//var closest *Triangle
	//lengh := 1000.

	for _, triangle := range scene.Triangles {
		if triangle == t {
			continue
		}
		if triangle != nil && !t.IsFollowedBy(triangle) {
			if t.Subscriber == nil {
				if t.GetColorDiff(triangle) < .1 && len(triangle.Followers) == 0 {
					t.Subscribe(triangle)
					triangle.FollowedBy(t)
				}
			} else {
				if t.GetColorDiff(triangle) < t.GetColorDiff(t.Subscriber) && len(triangle.Followers) == 0 {
					t.Subscriber.RemoveFollower(t)
					t.Subscribe(triangle)
					triangle.FollowedBy(t)
				}
			}
		}
	}

	if t.Subscriber != nil {
		if t.GetDistance(t.Subscriber) > 20 && t.GetDistance(t.Subscriber) < 200 {
			t.UpdateDirection(t.Subscriber.GetPos())
		}
		t.Speed = t.Subscriber.Speed
	}

	if t.G.Translation.X > scene.conf.MaxX {
		t.G.Translation.X = 0
		t.A.Translation.X = 0
		t.B.Translation.X = 0
		t.C.Translation.X = 0
	}
	if t.G.Translation.X < 0 {
		t.G.Translation.X = scene.conf.MaxX
		t.A.Translation.X = scene.conf.MaxX
		t.B.Translation.X = scene.conf.MaxX
		t.C.Translation.X = scene.conf.MaxX
	}

	if t.G.Translation.Y > scene.conf.MaxY {
		t.G.Translation.Y = 0
		t.A.Translation.Y = 0
		t.B.Translation.Y = 0
		t.C.Translation.Y = 0
	}
	if t.G.Translation.Y < 0 {
		t.G.Translation.Y = scene.conf.MaxY
		t.A.Translation.Y = scene.conf.MaxY
		t.B.Translation.Y = scene.conf.MaxY
		t.C.Translation.Y = scene.conf.MaxY
	}

}

func (t *Triangle) SameColor(triangle *Triangle) bool {
	diff := 0.3
	return math.Abs(t.Color.R-triangle.Color.R) < diff && math.Abs(t.Color.G-triangle.Color.G) < diff && math.Abs(t.Color.B-triangle.Color.B) < diff
}

func (t *Triangle) FollowedBy(triangle *Triangle) {
	t.Followers = append(t.Followers, triangle)
	//t.Color = pixel.RGB(1, 1, 1)
}

func (t *Triangle) Subscribe(triangle *Triangle) {
	t.Subscriber = triangle
}
func (t *Triangle) UnSubscribe() {
	t.Subscriber = nil
}

func (t *Triangle) IsFollowedBy(triangle *Triangle) bool {
	for _, follower := range t.Followers {
		if follower == triangle {
			return true
		}
		return follower.IsFollowedBy(triangle)
	}
	return false
}

func (t *Triangle) RemoveFollower(triangle *Triangle) {
	followers := make([]*Triangle, 0)
	for _, follower := range t.Followers {
		if follower != triangle {
			followers = append(followers, triangle)
		}
	}
	t.Followers = followers
}

func (t *Triangle) ApplyTranslation() {
	t.Translate(t.Direction.Scaled(t.Speed))
}

func (t *Triangle) Translate(vec pixel.Vec) {
	t.A.Translation = t.A.Translation.Add(vec)
	t.B.Translation = t.B.Translation.Add(vec)
	t.C.Translation = t.C.Translation.Add(vec)
	t.G.Translation = t.G.Translation.Add(vec)
}

func (t *Triangle) RefreshCenter() {
	t.G.Position.X = (t.A.Position.X + t.B.Position.X + t.C.Position.X) / 3
	t.G.Position.Y = (t.A.Position.Y + t.B.Position.Y + t.C.Position.Y) / 3
}

func (t *Triangle) ApplyRotation() {

	//t.DirectionInitial.ScaledXY(t.Direction).Scaled(-1 / (t.DirectionInitial.Len() * t.Direction.Len()))

	teta := t.Direction.Rotated(-t.DirectionInitial.Angle()).Angle()
	t.Rotate(teta)

	t.DirectionInitial = t.Direction
}

func (t *Triangle) updateTriangleSprite() {
	imd := imdraw.New(nil)
	t.imd = imd

	t.imd.Color = t.Color
	t.imd.Push(t.A.Position.Add(t.A.Translation))
	t.imd.Color = t.Color
	t.imd.Push(t.B.Position.Add(t.B.Translation))
	t.imd.Color = t.Color
	t.imd.Push(t.C.Position.Add(t.C.Translation))
	t.imd.Polygon(0)

	//t.imd.Color = pixel.RGB(0, 0, 1)
	//t.imd.Push(t.GetPos())
	//t.imd.Push(t.GetPos().Add(t.DirectionInitial.Scaled(30)))
	//
	//t.imd.Polygon(1)

	//if t.Subscriber != nil {
	//	t.imd.Color = pixel.RGB(0, 1, 0)
	//	t.imd.Push(t.GetPos())
	//	t.imd.Push(t.Subscriber.GetPos())
	//
	//	t.imd.Polygon(1)
	//}

}

func (t *Triangle) UpdateDirection(vec pixel.Vec) {
	t.Direction = vec.Sub(t.G.Position.Add(t.G.Translation)).Unit()
}
