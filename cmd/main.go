package main

import (
	"CHIP-8/system"
	e "github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const scale = 16

var chip8 *system.CHIP8

//func run(screen *e.Image) error {
//	chip8.Cycle()
//
//	var x = 0
//	var y = 0
//
//	if chip8.DrawFlag {
//		for i := 0; i < len(chip8.Gfx); i++ {
//			if i%64 == 0 && i >= 64{
//			//	fmt.Print(i)
//				x = 0
//				y += 1
//		//		fmt.Println()
//			}
//			if chip8.Gfx[i] == 1 {
//				screen.Set(x, y, color.White)
//			}
//			x++
//
//			//fmt.Print(chip8.Gfx[i])
//		}
//
//		//chip8.DrawFlag = false
//	}
//	return nil
//}

func run(screen *e.Image) error {
	chip8.Cycle()

	if chip8.DrawFlag {
		for x := 0; x < system.Width; x++ {
			for y := 0; y < system.Height; y++ {
				if chip8.Gfx[y][x] == 1 {
					screen.Set(x, y, color.White)
				}
			}
		}
	}
	return nil
}

func main() {
	chip8 = new(system.CHIP8)
	chip8.Initialize()
	err := chip8.Load("roms/pong.ch8")
	if err != nil {
		log.Fatal(err)
	}

	err = e.Run(run, system.Width, system.Height, scale, "CHIP-8")

	if err != nil {
		log.Fatal(err)
	}
}
