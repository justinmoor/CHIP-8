package main

import (
	"CHIP-8/chip8"
	"fmt"
	"syscall/js"
)

var (
	window, doc, body, canvas, ctx, gameOpts js.Value
	width                                    float64 = 64
	height                                   float64 = 32
	c                                        *chip8.CHIP8
	reset                                    chan string
)

const scale = 16

func setupHTML() {
	window = js.Global()
	doc = window.Get("document")
	body = doc.Get("body")

	gameOpts = doc.Call("createElement", "select")
	doc.Set("onchange", gameChange)
	games := [...]string{"PONG", "INVADERS"}

	for _, e := range games {
		game := doc.Call("createElement", "option")
		game.Set("value", e)
		game.Set("textContent", e)
		gameOpts.Call("appendChild", game)
	}

	br := doc.Call("createElement", "br")

	height = window.Get("innerHeight").Float()
	width = window.Get("innerWidth").Float()

	canvas = doc.Call("createElement", "canvas")
	canvas.Set("height", height)
	canvas.Set("width", width)

	ctx = canvas.Call("getContext", "2d")
	ctx.Set("fillStyle", "black")
	ctx.Call("scale", scale, scale)

	body.Call("appendChild", gameOpts)
	body.Call("appendChild", br)
	body.Call("appendChild", br)
	body.Call("appendChild", canvas)

	doc.Call("addEventListener", "keydown", keyEvent)
}

var gameChange = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	e := args[0]
	reset <- e.Get("target").Get("value").String()
	return nil
})

var keyEvent = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	e := args[0]
	keyPress := e.Get("key")

	if key, ok := chip8.KeyMap[keyPress.String()]; ok {
		c.SendKeyPress(key)
	}

	return nil
})

func main() {
	//quit := make(chan struct{})
	reset = make(chan string)

	setupHTML()

	gfxChan := make(chan [chip8.Width][chip8.Height]byte, 10)
	var gfxBuffer [chip8.Width][chip8.Height]byte

	LoadAndStartEmu("PONG", gfxChan)

	go func() {
		// web renderer
		var renderer js.Func
		renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			ctx.Call("clearRect", 0, 0, width, height)

			select {
			case gfxBuffer = <-gfxChan:
			default:
			}

			for x := 0; x < chip8.Width; x++ {
				for y := 0; y < chip8.Height; y++ {
					if gfxBuffer[x][y] == 1 {
						ctx.Call("fillRect", x, y, 1, 1)
					}
				}
			}

			window.Call("requestAnimationFrame", renderer)
			return nil
		})

		window.Call("requestAnimationFrame", renderer)

	}()

	select {}
}

func LoadAndStartEmu(game string, gfxChan chan [chip8.Width][chip8.Height]byte) {
	c = chip8.New()
	c.Logging = false

	// load a ROM from the static files
	if err := c.LoadRomHTTP(fmt.Sprintf("http://localhost:8000/roms/%v", game)); err != nil {
		panic(err)
	}

	// start the emulator
	go func() {
		for range c.Timer.C {
			c.Cycle()
			if c.DrawFlag {
				select {
				case gfxChan <- c.Gfx:
				case newGame := <-reset:
					c.Reset()
					if err := c.LoadRomHTTP(fmt.Sprintf("http://localhost:8000/roms/%v", newGame)); err != nil {
						panic(err)
					}
				default:
				}
				c.DrawFlag = false
			}
		}
	}()
}
