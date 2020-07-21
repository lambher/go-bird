package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	maxX = 500
	maxY = 500
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, maxX, maxY),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var scene Scene

	scene.GenerateScene(maxX, maxY)

	for !win.Closed() {
		win.Clear(colornames.Black)
		scene.CatchEvent(win)
		scene.Draw(win)
		scene.Update()
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
