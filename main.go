package main

import (
	"fmt"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  int32 = 1280
	windowHeight int32 = 720

	cameraFov  float32 = 45
	cameraNear float32 = 0.1
	cameraFar  float32 = 100.0
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

	window, err := sdl.CreateWindow("Learning Project", 200, 200, windowWidth, windowHeight, sdl.WINDOW_OPENGL)
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
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	cubePositions := []mgl32.Vec3{
		{-1.0, 0.0, 0.0},
		{5.0, 3.0, -10.0},
	}

	GenBindBuffer(gl.ARRAY_BUFFER) //VBO
	VAO := GenBindVertexArray()
	BufferData(gl.ARRAY_BUFFER, verticies, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, uintptr(3*4))
	gl.EnableVertexAttribArray(1)
	gl.BindVertexArray(0)

	keyboardState := sdl.GetKeyboardState()

	var camX, camY, camZ float32 = 0.0, 0.0, 0.0

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		if keyboardState[sdl.SCANCODE_W] != 0 {
			camZ += 0.1
		}
		if keyboardState[sdl.SCANCODE_S] != 0 {
			camZ -= 0.1
		}
		if keyboardState[sdl.SCANCODE_A] != 0 {
			camX += 0.1
		}
		if keyboardState[sdl.SCANCODE_D] != 0 {
			camX -= 0.1
		}
		if keyboardState[sdl.SCANCODE_SPACE] != 0 {
			camY -= 0.1
		}
		if keyboardState[sdl.SCANCODE_LSHIFT] != 0 {
			camY += 0.1
		}

		shaderProgram.Use()

		projMat := mgl32.Perspective(mgl32.DegToRad(cameraFov), float32(windowWidth)/float32(windowHeight), cameraNear, cameraFar)
		viewMat := mgl32.Ident4()
		viewMat = mgl32.Translate3D(camX, camY, camZ).Mul4(viewMat)
		shaderProgram.SetMatrix4("proj", projMat)
		shaderProgram.SetMatrix4("view", viewMat)

		BindTexture(texture)
		BindVertexArray(VAO)

		for i, pos := range cubePositions {
			modelMat := mgl32.Ident4()
			modelMat = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMat)
			angle := 25.0 * float32(i)
			modelMat = mgl32.HomogRotate3DY(mgl32.DegToRad(angle)).Mul4(modelMat)
			shaderProgram.SetMatrix4("model", modelMat)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.GLSwap()
		shaderProgram.CheckShadersForChanges()
	}
}
