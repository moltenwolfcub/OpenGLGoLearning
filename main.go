package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/moltenwolfcub/OpenGLGoLearning/helpers"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	cameraFov  float32 = 45
	cameraNear float32 = 0.1
	cameraFar  float32 = 100.0
)

var (
	windowWidth  int32 = 1280
	windowHeight int32 = 720
)

func main() {
	window, cleanup := helpers.SetupFPSWindow("Learning Project", windowWidth, windowHeight)
	defer cleanup()

	fmt.Println("OpenGL Version", helpers.GetVersion())

	window.WarpMouseInWindow(windowWidth/2, windowHeight/2)

	shaderProgram := helpers.NewShader("assets/shaders/test.vert", "assets/shaders/quadTexture.frag")
	texture := helpers.LoadTexture("assets/textures/metal/metalbox_full.png")

	cube := helpers.Cube(1)
	cubeBig := helpers.Cube(4)
	pent := helpers.Pentahedron(2)

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

	gl.BindVertexArray(0)

	focused := true
	focusCooldown := 0

	keyboardState := sdl.GetKeyboardState()

	camPos := mgl32.Vec3{0.0, 0.0, -2.0}
	worldUp := mgl32.Vec3{0.0, 1.0, 0.0}
	camera := helpers.NewCamera(camPos, worldUp, 90, 0, 0.0025, 0.1)

	elapsedTime := float32(0)
	for {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_RESIZED {
					windowWidth, windowHeight = e.Data1, e.Data2
					gl.Viewport(0, 0, windowWidth, windowHeight)
				}
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
			dirs := helpers.NewMoveDirs(
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

		helpers.BindTexture(texture)

		cube.DrawMultiple(shaderProgram, len(cubePositions), func(i int) mgl32.Mat4 {
			pos := cubePositions[i]
			return mgl32.Ident4().Mul4(mgl32.Translate3D(pos.X(), pos.Y(), pos.Z()))
		})
		cubeBig.Draw(shaderProgram, mgl32.Ident4().Mul4(mgl32.Translate3D(0, 5, 0)))
		pent.Draw(shaderProgram, mgl32.Ident4().Mul4(mgl32.Translate3D(0, 1, 0)))

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
