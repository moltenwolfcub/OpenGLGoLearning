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
