package main

import (
	"syscall/js"
	"unsafe"

	"github.com/bobcob7/wasm-rotating-cube/gltypes"
	"github.com/go-gl/mathgl/mgl32"
)

var (
	width   int
	height  int
	gl      js.Value
	glTypes gltypes.GLTypes
)

//// BUFFERS + SHADERS ////
// Shamelessly copied from https://www.tutorialspoint.com/webgl/webgl_cube_rotation.htm //
var verticesNative = []float32{
	-1, -1, -1, 1, -1, -1, 1, 1, -1, -1, 1, -1,
	-1, -1, 1, 1, -1, 1, 1, 1, 1, -1, 1, 1,
	-1, -1, -1, -1, 1, -1, -1, 1, 1, -1, -1, 1,
	1, -1, -1, 1, 1, -1, 1, 1, 1, 1, -1, 1,
	-1, -1, -1, -1, -1, 1, 1, -1, 1, 1, -1, -1,
	-1, 1, -1, -1, 1, 1, 1, 1, 1, 1, 1, -1,
}
var colorsNative = []float32{
	5, 3, 7, 5, 3, 7, 5, 3, 7, 5, 3, 7,
	1, 1, 3, 1, 1, 3, 1, 1, 3, 1, 1, 3,
	0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1,
	1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 0,
	1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 1, 0,
	0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0,
}
var indicesNative = []uint16{
	0, 1, 2, 0, 2, 3, 4, 5, 6, 4, 6, 7,
	8, 9, 10, 8, 10, 11, 12, 13, 14, 12, 14, 15,
	16, 17, 18, 16, 18, 19, 20, 21, 22, 20, 22, 23,
}

const vertShaderCode = `
attribute vec3 position;
uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;
attribute vec3 color;
varying vec3 vColor;

void main(void) {
	gl_Position = Pmatrix*Vmatrix*Mmatrix*vec4(position, 1.);
	vColor = color;
}
`
const fragShaderCode = `
precision mediump float;
varying vec3 vColor;
void main(void) {
	gl_FragColor = vec4(vColor, 1.);
}
`

func main() {
	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "gocanvas")
	width = doc.Get("body").Get("clientWidth").Int()
	height = doc.Get("body").Get("clientHeight").Int()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)

	gl = canvasEl.Call("getContext", "webgl")
	if gl.IsUndefined() {
		gl = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if gl.IsUndefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return
	}

	// Get some WebGL bindings
	glTypes.New(gl)

	// Convert buffers to JS TypedArrays
	var colors = gltypes.SliceToTypedArray(colorsNative)
	var vertices = gltypes.SliceToTypedArray(verticesNative)
	var indices = gltypes.SliceToTypedArray(indicesNative)

	// Create vertex buffer
	vertexBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", glTypes.ArrayBuffer, vertexBuffer)
	gl.Call("bufferData", glTypes.ArrayBuffer, vertices, glTypes.StaticDraw)

	// Create color buffer
	colorBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", glTypes.ArrayBuffer, colorBuffer)
	gl.Call("bufferData", glTypes.ArrayBuffer, colors, glTypes.StaticDraw)

	// Create index buffer
	indexBuffer := gl.Call("createBuffer")
	gl.Call("bindBuffer", glTypes.ElementArrayBuffer, indexBuffer)
	gl.Call("bufferData", glTypes.ElementArrayBuffer, indices, glTypes.StaticDraw)

	//// Shaders ////

	// Create a vertex shader object
	vertShader := gl.Call("createShader", glTypes.VertexShader)
	gl.Call("shaderSource", vertShader, vertShaderCode)
	gl.Call("compileShader", vertShader)

	// Create fragment shader object
	fragShader := gl.Call("createShader", glTypes.FragmentShader)
	gl.Call("shaderSource", fragShader, fragShaderCode)
	gl.Call("compileShader", fragShader)

	// Create a shader program object to store
	// the combined shader program
	shaderProgram := gl.Call("createProgram")
	gl.Call("attachShader", shaderProgram, vertShader)
	gl.Call("attachShader", shaderProgram, fragShader)
	gl.Call("linkProgram", shaderProgram)

	// Associate attributes to vertex shader
	PositionMatrix := gl.Call("getUniformLocation", shaderProgram, "Pmatrix")
	ViewMatrix := gl.Call("getUniformLocation", shaderProgram, "Vmatrix")
	ModelMatrix := gl.Call("getUniformLocation", shaderProgram, "Mmatrix")

	gl.Call("bindBuffer", glTypes.ArrayBuffer, vertexBuffer)
	position := gl.Call("getAttribLocation", shaderProgram, "position")
	gl.Call("vertexAttribPointer", position, 3, glTypes.Float, false, 0, 0)
	gl.Call("enableVertexAttribArray", position)

	gl.Call("bindBuffer", glTypes.ArrayBuffer, colorBuffer)
	color := gl.Call("getAttribLocation", shaderProgram, "color")
	gl.Call("vertexAttribPointer", color, 3, glTypes.Float, false, 0, 0)
	gl.Call("enableVertexAttribArray", color)

	gl.Call("useProgram", shaderProgram)

	// Set WeebGL properties
	gl.Call("clearColor", 0.5, 0.5, 0.5, 0.9) // Color the screen is cleared to
	gl.Call("clearDepth", 1.0)                // Z value that is set to the Depth buffer every frame
	gl.Call("viewport", 0, 0, width, height)  // Viewport size
	gl.Call("depthFunc", glTypes.LEqual)

	//// Create Matrixes ////
	ratio := float32(width) / float32(height)

	// Generate and apply projection matrix
	projMatrix := mgl32.Perspective(mgl32.DegToRad(45.0), ratio, 1, 100.0)
	var projMatrixBuffer *[16]float32
	projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&projMatrix))
	typedProjMatrixBuffer := gltypes.SliceToTypedArray([]float32((*projMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", PositionMatrix, false, typedProjMatrixBuffer)

	// Generate and apply view matrix
	viewMatrix := mgl32.LookAtV(mgl32.Vec3{3.0, 3.0, 3.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 1.0, 0.0})
	var viewMatrixBuffer *[16]float32
	viewMatrixBuffer = (*[16]float32)(unsafe.Pointer(&viewMatrix))
	typedViewMatrixBuffer := gltypes.SliceToTypedArray([]float32((*viewMatrixBuffer)[:]))
	gl.Call("uniformMatrix4fv", ViewMatrix, false, typedViewMatrixBuffer)

	//// Drawing the Cube ////
	movMatrix := mgl32.Ident4()
	var renderFrame js.Func
	var tmark float32
	var rotation = float32(0)

	// Bind to element array for draw function
	gl.Call("bindBuffer", glTypes.ElementArrayBuffer, indexBuffer)

	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// Calculate rotation rate
		now := float32(args[0].Float())
		tdiff := now - tmark
		tmark = now
		rotation = rotation + float32(tdiff)/500

		// Do new model matrix calculations
		movMatrix = mgl32.HomogRotate3DX(0.5 * rotation)
		movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DY(0.3 * rotation))
		movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DZ(0.2 * rotation))

		// Convert model matrix to a JS TypedArray
		var modelMatrixBuffer *[16]float32
		modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&movMatrix))
		typedModelMatrixBuffer := gltypes.SliceToTypedArray([]float32((*modelMatrixBuffer)[:]))

		// Apply the model matrix
		gl.Call("uniformMatrix4fv", ModelMatrix, false, typedModelMatrixBuffer)

		// Clear the screen
		gl.Call("enable", glTypes.DepthTest)
		gl.Call("clear", glTypes.ColorBufferBit)
		gl.Call("clear", glTypes.DepthBufferBit)

		// Draw the cube
		gl.Call("drawElements", glTypes.Triangles, len(indicesNative), glTypes.UnsignedShort, 0)

		// Call next frame
		js.Global().Call("requestAnimationFrame", renderFrame)

		return nil
	})
	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	done := make(chan struct{}, 0)
	<-done
}
