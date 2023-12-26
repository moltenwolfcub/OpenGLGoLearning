package helpers

import (
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

// used for initialising the target with the go array at data
// a more go-esque wrapper for the gl.BufferData function
func BufferData[T any](target uint32, data []T, usage uint32) {
	var v T
	dataTypeSize := unsafe.Sizeof(v)

	gl.BufferData(target, len(data)*int(dataTypeSize), gl.Ptr(data), usage)
}

type BufferID uint32

// used for generating and binding general buffers
// e.g. VertexBufferObject or NormalArrayObject
func GenBindBuffer(target uint32) BufferID {
	var buffer uint32
	gl.GenBuffers(1, &buffer)
	gl.BindBuffer(target, buffer)
	return BufferID(buffer)
}

// used for generating and binding vertex buffers
func GenBindVertexArray() BufferID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return BufferID(VAO)
}

func BindVertexArray(id BufferID) {
	gl.BindVertexArray(uint32(id))
}
