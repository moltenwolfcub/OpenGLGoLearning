package main

import (
	"fmt"
	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	window, err := sdl.CreateWindow("Learning Project", 200, 200, 1280, 720, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	fmt.Println("OpenGL Version", GetVersion())

	shaderProgram := NewShader("assets/shaders/test.vert", "assets/shaders/quadTexture.frag")
	texture := LoadTexture("assets/textures/test.png")

	//XYZ,UV
	verticies := []float32{
		0.5, 0.5, 0.0, 1.0, 1.0,
		0.5, -0.5, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.0, 0.0, 0.0,
		-0.5, 0.5, 0.0, 0.0, 1.0,
	}
	indices := []uint32{
		0, 1, 3, // triangle1
		1, 2, 3, // triangle2
	}

	GenBindBuffer(gl.ARRAY_BUFFER) //VBO
	VAO := GenBindVertexArray()
	BufferData(gl.ARRAY_BUFFER, verticies, gl.STATIC_DRAW)
	GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER) //EBO
	BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, uintptr(3*4))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	var x, y float32 = 0.0, 0.0

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		shaderProgram.Use()
		shaderProgram.SetFloat("x", float32(math.Sin(float64(x))))
		shaderProgram.SetFloat("y", float32(math.Cos(float64(y))))
		BindTexture(texture)
		BindVertexArray(VAO)
		gl.DrawElementsWithOffset(gl.TRIANGLES, 6, gl.UNSIGNED_INT, uintptr(0))

		window.GLSwap()
		shaderProgram.CheckShadersForChanges()
		x += 0.01
		y += 0.01
	}
}
