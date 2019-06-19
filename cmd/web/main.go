package main

import (
	"CHIP-8/system"
	"syscall/js"
)

var (
	window, doc, body, canvas, ctx, beep js.Value
	width                                float64 = 64
	height                               float64 = 32
)

const scale = 4

func setup() {
	window = js.Global()
	doc = window.Get("document")
	body = doc.Get("body")

	height = window.Get("innerHeight").Float()
	width = window.Get("innerWidth").Float()

	canvas = doc.Call("createElement", "canvas")
	canvas.Set("height", height)
	canvas.Set("width", width)
	body.Call("appendChild", canvas)

	ctx = canvas.Call("getContext", "2d")
	ctx.Set("fillStyle", "black")
	ctx.Call("scale", scale, scale)
}

func main() {
	setup()
	chip8 := system.New()

	if err := chip8.LoadRomHttp("http://localhost:8000/PONG1"); err != nil {
		panic(err)
	}

	var renderer js.Func
	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("clearRect", 0, 0, width, height)

		chip8.Cycle()
		if chip8.DrawFlag {
			for x := 0; x < system.Width; x++ {
				for y := 0; y < system.Height; y++ {
					if chip8.Gfx[y][x] == 1 {
						ctx.Call("fillRect", x, y, 1, 1)
					}
				}
			}
		}

		window.Call("requestAnimationFrame", renderer)
		return nil
	})

	window.Call("requestAnimationFrame", renderer)

	select {}
}
