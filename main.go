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

	cube := Cube(1)
	cubeBig := Cube(4)
	pent := Pentahedron(2)

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

type Object struct {
	verticies    []float32 //in XYZ UV
	vertexStride int       // 5 if using XYZ UV
	normals      []float32
	bufferLoader *helpers.BufferLoader
	vao          helpers.BufferID
	nao          helpers.BufferID
}

func (o *Object) fillBuffers() {
	o.bufferLoader = helpers.NewBufferLoader()
	o.vao = helpers.GenBindVertexArray()
	o.nao = helpers.GenBindBuffer(gl.ARRAY_BUFFER)

	helpers.GenBindBuffer(gl.ARRAY_BUFFER) //VBO

	helpers.BindVertexArray(o.vao)
	o.bufferLoader.BuildFloatBuffer(o.vao, helpers.NewBufferLayout([]int32{3, 2}, o.verticies))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(o.nao))
	o.bufferLoader.BuildFloatBuffer(o.nao, helpers.NewBufferLayout([]int32{3}, o.normals))
}

func (o *Object) calcNormals(triangleCount int) {
	vertexCount := triangleCount * 3 //3 bc we are working in 3d space so XYZ

	o.normals = make([]float32, vertexCount*3)
	for tri := 0; tri < triangleCount; tri++ {
		index := tri * o.vertexStride * 3
		p1 := mgl32.Vec3{o.verticies[index], o.verticies[index+1], o.verticies[index+2]}
		index += o.vertexStride
		p2 := mgl32.Vec3{o.verticies[index], o.verticies[index+1], o.verticies[index+2]}
		index += o.vertexStride
		p3 := mgl32.Vec3{o.verticies[index], o.verticies[index+1], o.verticies[index+2]}

		normal := helpers.TriangleNormal(p1, p2, p3)
		o.normals[tri*9+0] = normal.X()
		o.normals[tri*9+1] = normal.Y()
		o.normals[tri*9+2] = normal.Z()

		o.normals[tri*9+3] = normal.X()
		o.normals[tri*9+4] = normal.Y()
		o.normals[tri*9+5] = normal.Z()

		o.normals[tri*9+6] = normal.X()
		o.normals[tri*9+7] = normal.Y()
		o.normals[tri*9+8] = normal.Z()
	}
}

func (o Object) Draw(shader *helpers.Shader, drawMatrix mgl32.Mat4) {
	helpers.BindVertexArray(o.vao)

	shader.SetMatrix4("model", drawMatrix)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.verticies)/o.vertexStride))
}

func (o Object) DrawMultiple(shader *helpers.Shader, num int, drawMatrix func(int) mgl32.Mat4) {
	helpers.BindVertexArray(o.vao)

	for i := 0; i < num; i++ {
		shader.SetMatrix4("model", drawMatrix(i))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(o.verticies)/o.vertexStride))
	}
}

func Cube(size float32) Object {
	o := Object{}
	o.verticies = []float32{
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,
		-size / 2, size / 2, -size / 2, 0.0, 1.0,

		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, size / 2, size / 2, 1.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 1.0,
		-size / 2, size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,

		-size / 2, size / 2, size / 2, 1.0, 0.0,
		-size / 2, size / 2, -size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		-size / 2, size / 2, size / 2, 1.0, 0.0,

		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 0.0,

		-size / 2, -size / 2, -size / 2, 0.0, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 1.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
		-size / 2, -size / 2, -size / 2, 0.0, 1.0,

		-size / 2, size / 2, -size / 2, 0.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		size / 2, size / 2, -size / 2, 1.0, 1.0,
		size / 2, size / 2, size / 2, 1.0, 0.0,
		-size / 2, size / 2, -size / 2, 0.0, 1.0,
		-size / 2, size / 2, size / 2, 0.0, 0.0,
	}
	o.vertexStride = 5

	o.calcNormals(12)
	o.fillBuffers()

	return o
}

func Pentahedron(size float32) Object {
	o := Object{}
	o.verticies = []float32{
		size / 2, -size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 1.0,
		-size / 2, -size / 2, size / 2, 1.0, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		size / 2, -size / 2, -size / 2, 1.0, 0.0,
		-size / 2, -size / 2, -size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, -size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		-size / 2, -size / 2, size / 2, 1.0, 0.0,
		size / 2, -size / 2, size / 2, 0.0, 0.0,

		0.0, size / 2, 0.0, 0.5, 1.0,
		-size / 2, -size / 2, -size / 2, 1.0, 0.0,
		-size / 2, -size / 2, size / 2, 0.0, 0.0,
	}
	o.vertexStride = 5

	o.calcNormals(6)
	o.fillBuffers()

	return o
}
