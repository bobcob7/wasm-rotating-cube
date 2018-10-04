package gltypes

import "syscall/js"

type GLTypes struct {
	staticDraw         js.Value
	arrayBuffer        js.Value
	elementArrayBuffer js.Value
	vertexShader       js.Value
	fragmentShader     js.Value
	float              js.Value
	depthTest          js.Value
	colorBufferBit     js.Value
	depthBufferBit     js.Value
	triangles          js.Value
	unsignedShort      js.Value
	lEqual             js.Value
	lineLoop           js.Value
}

func (types *GLTypes) New() {
	types.staticDraw = gl.Get("STATIC_DRAW")
	types.arrayBuffer = gl.Get("ARRAY_BUFFER")
	types.elementArrayBuffer = gl.Get("ELEMENT_ARRAY_BUFFER")
	types.vertexShader = gl.Get("VERTEX_SHADER")
	types.fragmentShader = gl.Get("FRAGMENT_SHADER")
	types.float = gl.Get("FLOAT")
	types.depthTest = gl.Get("DEPTH_TEST")
	types.colorBufferBit = gl.Get("COLOR_BUFFER_BIT")
	types.triangles = gl.Get("TRIANGLES")
	types.unsignedShort = gl.Get("UNSIGNED_SHORT")
	types.lEqual = gl.Get("LEQUAL")
	types.depthBufferBit = gl.Get("DEPTH_BUFFER_BIT")
	types.lineLoop = gl.Get("LINE_LOOP")
}
