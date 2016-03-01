package main

import (
	"image"
	"image/png"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
)

func init() {
	runtime.LockOSThread()
}

func main() {

	glfwErr := glfw.Init()
	if glfwErr != nil {
		panic(glfwErr)
	}
	defer glfw.Terminate()

	monitor := glfw.GetPrimaryMonitor()
	vidMode := monitor.GetVideoMode()
	width := vidMode.Width
	height := vidMode.Height

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	window, err := glfw.CreateWindow(width, height, "readPixels", monitor, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	//window.SetSizeCallback(reshape)
	window.SetKeyCallback(onKey)
	window.SetCharCallback(onChar)
	glfw.SwapInterval(1)

	glErr := gl.Init()
	if glErr != nil {
		panic(glErr)
	}

	for !window.ShouldClose() {

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.ClearColor(0, 0, 0, 0)

		gl.LineWidth(1)

		gl.Begin(gl.QUADS)
		gl.Color3d(1, 1, 0)

		gl.Vertex2d(-.5, -.5)
		gl.Vertex2d(.5, -.5)
		gl.Vertex2d(.5, .5)
		gl.Vertex2d(-.5, .5)

		gl.End()

		gl.Flush()

		window.SwapBuffers()
		glfw.PollEvents()
		//		time.Sleep(2 * time.Second)

		saveFrame(width, height)
	}
}

func saveFrame(width int, height int) {

	screenshot := image.NewRGBA(image.Rect(0, 0, width, height))

	gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&screenshot.Pix))
	//gl.ReadPixels(0, 0, int32(width), int32(height), gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(&screenshot.Pix[0]))
	if gl.NO_ERROR != gl.GetError() {
		log.Println("panic pixels")
		panic("unable to read pixels")
	}

	filename := time.Now().Format("video/2006Jan02_15-04-05.999.png")

	os.Mkdir("video", os.ModeDir)
	outFile, _ := os.Create(filename)
	defer outFile.Close()

	png.Encode(outFile, screenshot)
}

func onChar(w *glfw.Window, char rune) {
	log.Println(char)
}

func onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	switch {
	case key == glfw.KeyEscape && action == glfw.Press,
		key == glfw.KeyQ && action == glfw.Press:
		w.SetShouldClose(true)
	}
}
