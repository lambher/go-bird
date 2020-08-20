package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/tkanos/gonfig"
	"golang.org/x/image/colornames"
)

type Conf struct {
	MaxX     float64
	MaxY     float64
	MinSpeed float64
	MaxSpeed float64
	Nb       int
}

func run() {
	var conf Conf

	err := gonfig.GetConf("./conf.json", &conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg := pixelgl.WindowConfig{
		Title:     "Pixel Rocks!",
		Bounds:    pixel.R(0, 0, conf.MaxX, conf.MaxY),
		VSync:     true,
		Resizable: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var scene Scene

	scene.GenerateScene(&conf)

	for !win.Closed() {
		win.Clear(colornames.Black)
		scene.CatchEvent(win)
		scene.Draw(win)
		scene.Update()
		win.Update()
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	pixelgl.Run(run)
}
