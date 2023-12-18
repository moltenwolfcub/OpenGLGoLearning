package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Learning Project", 200, 200, 1280, 720, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL Version", version)
}
