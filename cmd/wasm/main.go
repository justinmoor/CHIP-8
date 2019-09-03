package main

import (
	"CHIP-8/chip8"
	"fmt"
	"log"
	"syscall/js"
)

var (
	window, doc, body, canvas, ctx, gameOpts, beep js.Value
	width                                          float64 = 64
	height                                         float64 = 32
	c                                              *chip8.CHIP8
	reset                                          chan struct{}
	gfxChan                                        chan [chip8.Width][chip8.Height]byte
	game                                           string
)

const scale = 16

func setupHTML() {
	window = js.Global()
	doc = window.Get("document")
	body = doc.Get("body")

	// games options dropdown
	gameOpts = doc.Call("createElement", "select")
	doc.Set("onchange", gameChange)
	games := [...]string{"PONG", "PONG2", "INVADERS", "15PUZZLE", "BLINKY", "BLITZ", "BRIX", "CONNECT4", "GUESS", "HIDDEN", "KALEID",
		"MAZE", "MERLIN", "MISSILE", "PUZZLE", "SYZYGY", "TANK", "TETRIS", "TICTAC", "UFO", "VBRIX", "VERS", "WIPEOFF"}

	for _, e := range games {
		game := doc.Call("createElement", "option")
		game.Set("value", e)
		game.Set("textContent", e)
		gameOpts.Call("appendChild", game)
	}

	// add a break
	br := doc.Call("createElement", "br")

	// create canvas
	height = window.Get("innerHeight").Float()
	width = window.Get("innerWidth").Float()

	canvas = doc.Call("createElement", "canvas")
	canvas.Set("height", height)
	canvas.Set("width", width)

	ctx = canvas.Call("getContext", "2d")
	ctx.Set("fillStyle", "black")
	ctx.Call("scale", scale, scale)

	// add all elements to the document
	body.Call("appendChild", gameOpts)
	body.Call("appendChild", br)
	body.Call("appendChild", br)
	body.Call("appendChild", canvas)

	// add eventlistener
	doc.Call("addEventListener", "keydown", keyEvent)
}

var gameChange = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	e := args[0]
	game = e.Get("target").Get("value").String()
	close(reset)
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
	setupHTML()
	game = "PONG"

	for {
		reset = make(chan struct{})
		gfxChan = make(chan [chip8.Width][chip8.Height]byte, 10)
		var gfxBuffer [chip8.Width][chip8.Height]byte

		loadAndStartEmu(game, gfxChan)

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

		select {
		case <-reset:
			// reload everything
			continue
		}
	}

}

func loadAndStartEmu(game string, gfxChan chan [chip8.Width][chip8.Height]byte) {
	c = chip8.New()
	c.Logging = false

	// load a ROM from the static files
	if err := c.LoadRomHTTP(fmt.Sprintf("http://localhost:8000/roms/%v", game)); err != nil {
		panic(err)
	}

	// start the emulator
	go func() {
		for range c.Timer.C {
			select {
			case <-reset:
				return
			default:
				if err := c.Cycle(); err != nil {
					log.Fatal(err)
				}

				if c.DrawFlag {
					gfxChan <- c.Gfx
					c.DrawFlag = false
				}
			}
		}
	}()
}
