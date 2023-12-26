package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

type ShaderID uint32
type ProgramID uint32
type BufferID uint32
type TextureID uint32

func LoadTexture(filename string) TextureID {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[i] = byte(r / 256)
			i++
			pixels[i] = byte(g / 256)
			i++
			pixels[i] = byte(b / 256)
			i++
			pixels[i] = byte(a / 256)
			i++
		}
	}

	texture := GenBindTexture()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	gl.GenerateMipmap(gl.TEXTURE_2D)
	return texture
}

func GenBindTexture() TextureID {
	var id uint32
	gl.GenTextures(1, &id)
	gl.BindTexture(gl.TEXTURE_2D, id)
	return TextureID(id)
}
func BindTexture(id TextureID) {
	gl.BindTexture(gl.TEXTURE_2D, uint32(id))
}

func LoadShader(path string, shaderType uint32) ShaderID {
	shaderFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	shaderId := CreateShader(string(shaderFile), shaderType)
	return shaderId
}

func CreateShader(shaderSource string, shaderType uint32) ShaderID {
	shaderId := gl.CreateShader(shaderType)
	shaderSource += "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		panic("Failed to compile shader:\n" + log)
	}
	return ShaderID(shaderId)
}

func CreateProgram(vertPath string, fragPath string) ProgramID {
	vert := LoadShader(vertPath, gl.VERTEX_SHADER)
	frag := LoadShader(fragPath, gl.FRAGMENT_SHADER)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)
	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		panic("Failed to link program:\n" + log)
	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return ProgramID(shaderProgram)
}

func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return BufferID(buffer)
}
func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func GenEBO() BufferID {
	var EBO uint32
	gl.GenBuffers(1, &EBO)
	return BufferID(EBO)
}

func BufferData[T any](target uint32, data []T, usage uint32) {
	var v T
	dataTypeSize := unsafe.Sizeof(v)

	gl.BufferData(target, len(data)*int(dataTypeSize), gl.Ptr(data), usage)
}

func UseProgram(id ProgramID) {
	gl.UseProgram(uint32(id))
}

func BindVertexArray(id BufferID) {
	gl.BindVertexArray(uint32(id))
}

func TriangleNormal(p1, p2, p3 mgl32.Vec3) mgl32.Vec3 {
	U := p2.Sub(p1)
	V := p3.Sub(p1)

	return U.Cross(V).Normalize()
}

// shaders
type Shader struct {
	id          ProgramID
	vertPath    string
	vertModTime time.Time
	fragPath    string
	fragModTime time.Time
}

func NewShader(vertPath string, fragPath string) *Shader {
	id := CreateProgram(vertPath, fragPath)

	s := Shader{
		id:       id,
		vertPath: vertPath,
		fragPath: fragPath,

		vertModTime: getModTime(vertPath),
		fragModTime: getModTime(fragPath),
	}

	return &s
}

func (s *Shader) Use() {
	UseProgram(s.id)
}

func (s *Shader) SetFloat(name string, value float32) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	gl.Uniform1f(loc, value)
}
func (s *Shader) SetMatrix4(name string, value mgl32.Mat4) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	m4 := [16]float32(value)
	gl.UniformMatrix4fv(loc, 1, false, &m4[0])
}
func (s *Shader) SetVec3(name string, value mgl32.Vec3) {
	name_cstr := gl.Str(name + "\x00")
	loc := gl.GetUniformLocation(uint32(s.id), name_cstr)

	v3 := [3]float32(value)
	gl.Uniform3fv(loc, 1, &v3[0])
}

func (s *Shader) CheckShadersForChanges() {
	vertModTime := getModTime(s.vertPath)
	fragModTime := getModTime(s.fragPath)
	if v, f := !vertModTime.Equal(s.vertModTime), !fragModTime.Equal(s.fragModTime); v || f {
		if v {
			fmt.Printf("A vertex shader file has been modified: %s\n", s.vertPath)
			s.vertModTime = vertModTime
		}
		if f {
			fmt.Printf("A fragment shader file has been modified: %s\n", s.fragPath)
			s.fragModTime = fragModTime
		}
		id := CreateProgram(s.vertPath, s.fragPath)

		gl.DeleteProgram(uint32(s.id))
		s.id = id
	}
}

func getModTime(path string) time.Time {
	file, err := os.Stat(path)
	if err != nil {
		panic(err)
	}
	return file.ModTime()
}

//camera

type MovementDirs struct {
	Forward int
	Right   int
	Up      int
}

func NewMoveDirs(f, b, r, l, u, d bool) MovementDirs {
	var fi, bi, ri, li, ui, di int
	if f {
		fi = 1
	}
	if b {
		bi = 1
	}
	if r {
		ri = 1
	}
	if l {
		li = 1
	}
	if u {
		ui = 1
	}
	if d {
		di = 1
	}

	moveDirs := MovementDirs{
		Forward: fi - bi,
		Right:   ri - li,
		Up:      ui - di,
	}
	return moveDirs
}

type Camera struct {
	Pos mgl32.Vec3

	Up      mgl32.Vec3
	Right   mgl32.Vec3
	Forward mgl32.Vec3

	WorldUp mgl32.Vec3

	Yaw   float32
	Pitch float32

	MovementSpeed    float32
	MouseSensitivity float32
	Zoom             float32
}

func NewCamera(pos, worldUp mgl32.Vec3, yaw, pitch, speed, sensitivity float32) *Camera {
	cam := Camera{
		Pos:              pos,
		WorldUp:          worldUp,
		Yaw:              yaw,
		Pitch:            pitch,
		MovementSpeed:    speed,
		MouseSensitivity: sensitivity,
	}
	cam.updateVectors()

	return &cam
}

func (c *Camera) updateVectors() {
	forward := mgl32.Vec3{
		Cos32Deg(c.Yaw) * Cos32Deg(c.Pitch),
		Sin32Deg(c.Pitch),
		Sin32Deg(c.Yaw) * Cos32Deg(c.Pitch),
	}

	c.Forward = forward.Normalize()
	c.Right = forward.Cross(c.WorldUp).Normalize()
	c.Up = c.Right.Cross(c.Forward).Normalize()
}

func (c *Camera) GetViewMatrix() mgl32.Mat4 {
	center := c.Pos.Add(c.Forward)

	return mgl32.LookAt(
		c.Pos.X(), c.Pos.Y(), c.Pos.Z(),
		center.X(), center.Y(), center.Z(),
		c.Up.X(), c.Up.Y(), c.Up.Z(),
	)
}

func (c *Camera) UpdateCamera(dir MovementDirs, deltaTime, mouseDx, mouseDy float32) {
	magnitude := c.MovementSpeed * deltaTime

	//remove Z component and normalize
	forwardMovement := mgl32.Vec3{c.Forward.X(), 0, c.Forward.Z()}
	if forwardMovement.Len() > 0 {
		forwardMovement = forwardMovement.Normalize()
	}

	c.Pos = c.Pos.Add(forwardMovement.Mul(magnitude).Mul(float32(dir.Forward)))
	c.Pos = c.Pos.Add(c.Right.Mul(magnitude).Mul(float32(dir.Right)))
	c.Pos = c.Pos.Add(c.WorldUp.Mul(magnitude).Mul(float32(dir.Up)))

	mouseDx *= c.MouseSensitivity
	mouseDy *= c.MouseSensitivity

	c.Yaw += mouseDx
	if c.Yaw < 0 {
		c.Yaw = 360 - mgl32.Abs(c.Yaw)
	} else if c.Yaw >= 360 {
		c.Yaw -= 360
	}

	c.Pitch += mouseDy
	if c.Pitch >= 90 {
		c.Pitch = 89.9999
	} else if c.Pitch <= -90 {
		c.Pitch = -89.9999
	}

	c.updateVectors()
}
