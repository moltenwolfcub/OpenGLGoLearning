package main

import (
	"fmt"
	"time"

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
	sdl.SetRelativeMouseMode(true)
	if err != nil {
		panic(err)
	}
	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.CULL_FACE)
	fmt.Println("OpenGL Version", GetVersion())

	window.WarpMouseInWindow(windowWidth/2, windowHeight/2)

	shaderProgram := NewShader("assets/shaders/test.vert", "assets/shaders/quadTexture.frag")
	texture := LoadTexture("assets/textures/metal/metalbox_full.png")

	//XYZ,UV
	verticies := []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,

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
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
	}

	normals := make([]float32, 36*3)
	for tri := 0; tri < 12; tri++ {
		index := tri * 15
		p1 := mgl32.Vec3{verticies[index], verticies[index+1], verticies[index+2]}
		index += 5
		p2 := mgl32.Vec3{verticies[index], verticies[index+1], verticies[index+2]}
		index += 5
		p3 := mgl32.Vec3{verticies[index], verticies[index+1], verticies[index+2]}

		normal := TriangleNormal(p1, p2, p3)
		normals[tri*9+0] = normal.X()
		normals[tri*9+1] = normal.Y()
		normals[tri*9+2] = normal.Z()

		normals[tri*9+3] = normal.X()
		normals[tri*9+4] = normal.Y()
		normals[tri*9+5] = normal.Z()

		normals[tri*9+6] = normal.X()
		normals[tri*9+7] = normal.Y()
		normals[tri*9+8] = normal.Z()
	}

	cubePositions := []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{1.1, 0.0, 0.0},
		{2.2, 0.0, 0.0},
		{3.3, 0.0, 0.0},
		{4.4, 0.0, 0.0},
		{5.5, 0.0, 0.0},

		{5.0, 1.0, -5.0},
		{-5.0, -2.0, 1.0},
	}

	GenBindBuffer(gl.ARRAY_BUFFER) //VBO
	VAO := GenBindVertexArray()
	BufferData(gl.ARRAY_BUFFER, verticies, gl.STATIC_DRAW)

	BindVertexArray(VAO)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 5*4, uintptr(3*4))
	gl.EnableVertexAttribArray(1)

	NAO := GenBindBuffer(gl.ARRAY_BUFFER)
	BufferData(gl.ARRAY_BUFFER, normals, gl.STATIC_DRAW)

	BindVertexArray(NAO)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(2)

	gl.BindVertexArray(0)

	focused := true
	focusCooldown := 0

	keyboardState := sdl.GetKeyboardState()

	camPos := mgl32.Vec3{0.0, 0.0, -2.0}
	worldUp := mgl32.Vec3{0.0, 1.0, 0.0}
	camera := NewCamera(camPos, worldUp, 90, 0, 0.0025, 0.1)

	elapsedTime := float32(0)
	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		if keyboardState[sdl.SCANCODE_ESCAPE] != 0 {
			return
		}
		if keyboardState[sdl.SCANCODE_I] != 0 {
			fmt.Printf("Yaw: %v, Pitch %v\n", camera.Yaw, camera.Pitch)
		}
		if focusCooldown == 0 {
			if keyboardState[sdl.SCANCODE_F] != 0 {
				focused = !focused
				sdl.SetRelativeMouseMode(focused)

				if focused {
					window.WarpMouseInWindow(windowWidth/2, windowHeight/2)
				}

				focusCooldown = 10
			}
		} else {
			if focusCooldown < 0 {
				focusCooldown = 0
			} else {
				focusCooldown--
			}
		}

		gl.ClearColor(0.0, 0.0, 0.0, 0.0)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		if focused {
			dirs := NewMoveDirs(
				keyboardState[sdl.SCANCODE_W] != 0,
				keyboardState[sdl.SCANCODE_S] != 0,
				keyboardState[sdl.SCANCODE_D] != 0,
				keyboardState[sdl.SCANCODE_A] != 0,
				keyboardState[sdl.SCANCODE_SPACE] != 0,
				keyboardState[sdl.SCANCODE_LSHIFT] != 0,
			)
			mouseX, mouseY, _ := sdl.GetMouseState()
			mouseDx, mouseDy := float32(mouseX-windowWidth/2), -float32(mouseY-windowHeight/2)

			camera.UpdateCamera(dirs, elapsedTime, mouseDx, mouseDy)
		}
		shaderProgram.Use()

		projMat := mgl32.Perspective(mgl32.DegToRad(cameraFov), float32(windowWidth)/float32(windowHeight), cameraNear, cameraFar)
		viewMat := camera.GetViewMatrix()
		shaderProgram.SetMatrix4("proj", projMat)
		shaderProgram.SetMatrix4("view", viewMat)

		shaderProgram.SetVec3("viewPos", camera.Pos)
		shaderProgram.SetVec3("lightPos", mgl32.Vec3{3.3, 1, 0})
		shaderProgram.SetVec3("lightColor", mgl32.Vec3{1, 1, 1})
		shaderProgram.SetVec3("ambientLight", mgl32.Vec3{0.3, 0.3, 0.3})

		BindTexture(texture)
		BindVertexArray(VAO)

		for i, pos := range cubePositions {
			modelMat := mgl32.Ident4()

			angle := 25.0 * float32(i) * 0
			modelMat = mgl32.HomogRotate3D(mgl32.DegToRad(angle), mgl32.Vec3{1, 0, 0}).Mul4(modelMat)

			modelMat = mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()).Mul4(modelMat)

			shaderProgram.SetMatrix4("model", modelMat)
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

		window.GLSwap()
		shaderProgram.CheckShadersForChanges()

		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)

		if focused {
			sdl.EventState(sdl.MOUSEMOTION, sdl.IGNORE)
			window.WarpMouseInWindow(windowWidth/2, windowHeight/2)
			sdl.EventState(sdl.MOUSEMOTION, sdl.ENABLE)
		}
	}
}
