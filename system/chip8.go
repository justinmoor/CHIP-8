package system

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
)

const (
	Width  = 64
	Height = 32
)

type CHIP8 struct {
	cpu
	ScreenState chan [Height][Width]byte
	KeyPress    chan uint16
	gfx         [Height][Width]byte // display
	keys        [16]byte            // current key state
	drawFlag    bool
	delayTimer  byte
	soundTimer  byte
}

type cpu struct {
	opcode uint16     // opcode is 2 bytes
	memory [4096]byte // RAM
	v      [16]byte   // CPU registers
	i      uint16     // index register
	pc     uint16     // program counter
	stack  [16]uint16 // stack
	sp     uint16     // stackpointer
}

func (c *CHIP8) Run(rom string) error {
	c.initialize()
	if err := c.load(rom); err != nil {
		return err
	}

	go func() {
		for {
			c.cycle()

			if c.drawFlag {
				select {
				case c.ScreenState <- c.gfx:
					break
				default:
					break
				}
			}
		}
	}()

	return nil
}

func (c *CHIP8) SendKeyPress(key uint16) {
	select {
	case c.KeyPress <- key:
		break
	default:
		break
	}
}

func (c *CHIP8) initialize() {
	c.pc = 0x200 // = 512: initial point where a program will start
	c.opcode = 0x00
	c.i = 0x00
	c.sp = 0x00
	c.drawFlag = false

	// clear all the memory
	c.memory = [4096]byte{}
	c.v = [16]byte{}
	c.stack = [16]uint16{}

	c.ScreenState = make(chan [Height][Width]byte)
	c.KeyPress = make(chan uint16)

	// load font into memory
	for i := 0; i < 0x50; i++ {
		c.memory[i] = font[i]
	}
}

// fetch opcode
// decode opcode
// execude opcode
// update timers
func (c *CHIP8) cycle() {
	c.opcode = uint16(c.memory[c.pc])
	c.opcode <<= 8
	c.opcode |= uint16(c.memory[c.pc+1])

	addr := c.opcode & 0x0FFF
	nn := byte(c.opcode & 0x00FF)
	n := byte(c.opcode & 0x000F)
	x := byte((c.opcode & 0x0F00) >> 8)
	y := byte((c.opcode & 0x00F0) >> 4)

	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode & 0x000F {
		case 0x0000:
			c.exec00E0()
			break
		case 0x000E:
			c.exec00EE()
			break
		default:
			fmt.Println("0000: Unsupported opcode")
			break
		}
		break
	case 0x1000:
		c.exec1NNN(addr)
		break
	case 0x2000:
		c.exec2NNN(addr)
		break
	case 0x3000:
		c.exec3XNN(x, nn)
		break
	case 0x4000:
		c.exec4XNN(x, nn)
		break
	case 0x5000:
		c.exec5XY0(x, y)
		break
	case 0x6000:
		c.exec6XNN(x, nn)
		break
	case 0x7000:
		c.exec7XNN(x, nn)
		break
	case 0x8000:
		switch c.opcode & 0x000F {
		case 0x0000:
			c.exec8XY0(x, y)
			break
		case 0x0001:
			c.exec8XY1(x, y)
			break
		case 0x0002:
			c.exec8XY2(x, y)
			break
		case 0x0003:
			c.exec8XY3(x, y)
			break
		case 0x0004:
			c.exec8XY4(x, y)
			break
		case 0x0005:
			c.exec8XY5(x, y)
			break
		case 0x0006:
			c.exec8XY6(x, y)
			break
		case 0x0007:
			c.exec8XY7(x, y)
			break
		case 0x000E:
			c.exec8XYE(x, y)
			break
		default:
			fmt.Println("8000: Unsupported opcode")
			break
		}
		break
	case 0x9000:
		c.exec9XY0(x, y)
		break
	case 0xA000:
		c.execANNN(addr)
		break
	case 0xB000:
		c.execBNNN(addr)
		break
	case 0xC000:
		c.execCXNN(x, nn)
		break
	case 0xD000:
		c.execDXYN(x, y, n)
		break
	case 0xE000:
		switch c.opcode & 0x00FF {
		case 0x009E:
			c.execEX9E(x)
			break
		case 0x00A1:
			c.execEXA1(x)
			break
		}
		break
	case 0xF000:
		switch c.opcode & 0x0FF {
		case 0x0007:
			c.execFX07(x)
			break
		case 0x000A:
			c.execFX0A(x)
			break
		case 0x0015:
			c.execFX15(x)
			break
		case 0x0018:
			c.execFX18(x)
			break
		case 0x001E:
			c.execFX1E(x)
			break
		case 0x0029:
			c.execFX29(x)
			break
		case 0x0033:
			c.execFX33(x)
			break
		case 0x0055:
			c.execFX55(x)
			break
		case 0x0065:
			c.execFX65(x)
			break
		default:
			//fmt.Println("FX00: Unsupported opcode")
			break
		}
		break
	case 0x0033:
		c.execFX33(x)
		break
	default:

	}

	//if c.DrawFlag {
	//	c.debugDraw()
	//}
	//fmt.Printf("%s\n", fmt.Sprintf("%x", c.opcode))
	//fmt.Printf("%s\n", fmt.Sprintf("%x", c.pc))
}

//func (c *CHIP8) debugDraw() {
//	for x := 0; x < len(c.Gfx); x++ {
//		if x%64 == 0 {
//			fmt.Println()
//		}
//
//		if c.Gfx[x] == 1 {
//			fmt.Print("1")
//		} else {
//			fmt.Print("0")
//		}
//	}
//
//	fmt.Println()
//}

func (c *CHIP8) exec00E0() {
	fmt.Println("Executing 00E0")
	c.gfx = [Height][Width]byte{}
	c.drawFlag = true
	c.pc += 2
}

func (c *CHIP8) exec00EE() {
	fmt.Println("Executing 00EE")
	c.sp--
	c.pc = c.stack[c.sp]
}

func (c *CHIP8) exec1NNN(addr uint16) {
	fmt.Println("Executing 1NNN")
	c.pc = addr
}

func (c *CHIP8) exec2NNN(addr uint16) {
	fmt.Println("Executing 2NNN")
	c.stack[c.sp] = c.pc + 2
	c.sp++
	c.pc = addr
}

func (c *CHIP8) exec3XNN(x, nn byte) {
	fmt.Println("Executing 3XNN")
	if c.v[x] == nn {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec4XNN(x, nn byte) {
	fmt.Println("Executing 5XY0")
	if c.v[x] != nn {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec5XY0(x, y byte) {
	fmt.Println("Executing 5XY0")
	if c.v[x] == c.v[y] {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec6XNN(x, nn byte) {
	fmt.Println("Executing 6XNN")
	c.v[x] = nn
	c.pc += 2
}

func (c *CHIP8) exec7XNN(x, nn byte) {
	fmt.Println("Executing 7XNN")
	c.v[x] += nn
	c.pc += 2
}

func (c *CHIP8) exec8XY0(x, y byte) {
	fmt.Println("Executing 8XY0")
	c.v[x] = c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY1(x, y byte) {
	fmt.Println("Executing 8XY1")
	c.v[x] |= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY2(x, y byte) {
	fmt.Println("Executing 8XY2")
	c.v[x] &= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY3(x, y byte) {
	fmt.Println("Executing 8XY3")
	c.v[x] ^= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY4(x, y byte) {
	fmt.Println("Executing 8XY4")
	if c.v[y] > (0xFF - c.v[y]) {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
	c.v[x] += c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY5(x, y byte) {
	fmt.Println("Executing 8XY5")
	if c.v[y] > (0xFF - c.v[x]) {
		c.v[0xF] = 0
	} else {
		c.v[0xF] = 1
	}
	c.v[x] -= c.v[y]
	c.pc += 2
}

//TODO: probably not correct
func (c *CHIP8) exec8XY6(x, y byte) {
	fmt.Println("Executing 8XY6")
	c.v[0xF] = c.v[x] & 0x0001
	c.v[x] >>= 1
	c.pc += 2
}

func (c *CHIP8) exec8XY7(x, y byte) {
	fmt.Println("Executing 8XY7")
	c.v[x] = c.v[y] - c.v[x]
	c.pc += 2
}

func (c *CHIP8) exec8XYE(x, y byte) {
	fmt.Println("Executing 8XYE")
	c.v[0xF] = (c.v[x] >> 7) & 0x0001
	c.v[x] <<= 1
	c.pc += 2
}

func (c *CHIP8) exec9XY0(x, y byte) {
	fmt.Println("Executing 9XY0")
	if c.v[x] != c.v[y] {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) execANNN(addr uint16) {
	fmt.Println("Executing ANNN")
	c.i = addr
	c.pc += 2
}

func (c *CHIP8) execBNNN(addr uint16) {
	fmt.Println("Executing BNNNN")
	c.pc = addr + uint16(c.v[0x0000])
}

func (c *CHIP8) execCXNN(x, nn byte) {
	fmt.Println("Executing CXNN")
	b := randomByte()
	c.v[x] = b & nn
	c.pc += 2
}

func (c *CHIP8) execDXYN(x, y, n byte) {
	fmt.Println("Executing DXYN")
	vx := c.v[x]
	vy := c.v[y]
	var pixel byte

	c.v[0xF] = 0
	for yl := byte(0); yl < n; yl++ { // n = height
		pixel = c.memory[c.i+uint16(yl)]
		for xl := byte(0); xl < 8; xl++ { // width => always 8 pixels

			if (pixel & (0x80 >> xl)) != 0 {
				if c.gfx[(vy+yl)%Height][(vx+xl)%Width] == 1 {
					c.v[0xF] = 1
				}
				c.gfx[(y + yl)][(vx + xl)] ^= 1
			}
		}
	}
	c.drawFlag = true
	c.pc += 2
}

func (c *CHIP8) execEX9E(x byte) {
	fmt.Println("Executing EX9E")
	if c.keys[c.v[x]] != 0 {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) execEXA1(x byte) {
	fmt.Println("Executing EXA1")
	if c.keys[c.v[x]] == 0 {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) execFX07(x byte) {
	fmt.Println("Executing FX07")
	c.v[x] = c.delayTimer
	c.pc += 2
}

func (c *CHIP8) execFX0A(x byte) {
	fmt.Println("Executing FX0A")
	c.v[x] = c.keys[<-c.KeyPress]
	c.pc += 2
}

func (c *CHIP8) execFX15(x byte) {
	fmt.Println("Executing FX15")
	c.delayTimer = c.v[x]
	c.pc += 2
}

func (c *CHIP8) execFX18(x byte) {
	fmt.Println("Executing FX18")
	c.soundTimer = c.v[x]
	c.pc += 2
}

func (c *CHIP8) execFX1E(x byte) {
	fmt.Println("Executing FX1E")
	c.i += uint16(c.v[x])
	c.pc += 2
}

func (c *CHIP8) execFX29(x byte) {
	fmt.Println("Executing FX29")
	c.i = uint16(font[c.v[x]])
	c.pc += 2
}

func (c *CHIP8) execFX33(x byte) {
	fmt.Println("Executing FX33")
	c.memory[c.i] = c.v[x] / 100
	c.memory[c.i+1] = (c.v[x] / 10) % 10
	c.memory[c.i+2] = (c.v[x] % 100) % 10
	c.pc += 2
}

func (c *CHIP8) execFX55(x byte) {
	fmt.Println("Executing FX55")
	for i := byte(0); i < x; i++ {
		c.memory[c.i+uint16(i)] = c.v[i]
	}
	c.i += uint16(x) + 1
	c.pc += 2
}

func (c *CHIP8) execFX65(x byte) {
	fmt.Println("Executing FX65")
	for i := byte(0); i < x; i++ {
		c.v[i] = c.memory[c.i+uint16(i)]
	}
	c.i += uint16(x) + 1
	c.pc += 2
}

func (c *CHIP8) load(romName string) error {
	rom, err := os.Open(romName)
	if err != nil {
		return err
	}
	defer rom.Close()

	stats, err := rom.Stat()
	if err != nil {
		return err
	}

	reader := bufio.NewReader(rom)
	buffer := make([]byte, stats.Size())
	if _, err = reader.Read(buffer); err != nil {
		return err
	}

	// read ROM into memory
	for i := 0; i < len(buffer); i++ {
		c.memory[i+0x200] = buffer[i]
		//hex := fmt.Sprintf("%x", c.memory[i+0x200])
		//fmt.Printf("%s\n", hex)
	}
	return nil
}

func randomByte() byte {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]
}
