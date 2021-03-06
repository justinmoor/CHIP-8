package chip8

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	Width  = 64
	Height = 32
)

var KeyMap = map[string]uint8{
	"1": 0x1, "2": 0x2, "3": 0x3, "4": 0xC,
	"q": 0x4, "w": 0x5, "e": 0x6, "r": 0xD,
	"a": 0x7, "s": 0x8, "d": 0x9, "f": 0xE,
	"z": 0xA, "x": 0x0, "c": 0xB, "v": 0xF,
}

type CHIP8 struct {
	cpu
	memory     [4096]byte          // RAM
	Gfx        [Width][Height]byte // display buffer
	keys       [16]byte            // current key state
	DrawFlag   bool
	delayTimer byte
	SoundTimer byte
	Timer      *time.Ticker
	Logging    bool
}

type cpu struct {
	opcode uint16     // opcode is 2 bytes
	v      [16]byte   // CPU registers
	i      uint16     // index register
	pc     uint16     // program counter
	stack  [16]uint16 // stack
	sp     uint16     // stackpointer
}

func New() *CHIP8 {
	c := new(CHIP8)
	c.Reset()
	return c
}

func (c *CHIP8) SendKeyPress(k uint8) {
	c.keys[k] = 1
}

func (c *CHIP8) Reset() {
	c.pc = 0x200 // = 512: initial point where a program will start
	c.opcode = 0x00
	c.i = 0x00
	c.sp = 0x00
	c.DrawFlag = false
	c.Timer = time.NewTicker(2 * time.Millisecond)

	// clear all the memory
	c.memory = [4096]byte{}
	c.v = [16]byte{}
	c.stack = [16]uint16{}

	// load font into memory
	for i := 0; i < 0x50; i++ {
		c.memory[i] = font[i]
	}
}

// fetch opcode
// decode opcode
// execude opcode
// update timers
func (c *CHIP8) Cycle() error {
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
		switch c.opcode & 0x00FF {
		case 0x0000:
			c.exec0NNN(addr)
		case 0x00E0:
			c.exec00E0()
		case 0x00EE:
			c.exec00EE()
		default:
			return errors.New(fmt.Sprintf("0000: Unsupported opcode %v", c.opcode))
		}
	case 0x1000:
		c.exec1NNN(addr)
	case 0x2000:
		c.exec2NNN(addr)
	case 0x3000:
		c.exec3XNN(x, nn)
	case 0x4000:
		c.exec4XNN(x, nn)
	case 0x5000:
		c.exec5XY0(x, y)
	case 0x6000:
		c.exec6XNN(x, nn)
	case 0x7000:
		c.exec7XNN(x, nn)
	case 0x8000:
		switch c.opcode & 0x000F {
		case 0x0000:
			c.exec8XY0(x, y)
		case 0x0001:
			c.exec8XY1(x, y)
		case 0x0002:
			c.exec8XY2(x, y)
		case 0x0003:
			c.exec8XY3(x, y)
		case 0x0004:
			c.exec8XY4(x, y)
		case 0x0005:
			c.exec8XY5(x, y)
		case 0x0006:
			c.exec8XY6(x, y)
		case 0x0007:
			c.exec8XY7(x, y)
		case 0x000E:
			c.exec8XYE(x, y)
		default:
			return errors.New(fmt.Sprintf("8000: Unsupported opcode %v", c.opcode))
		}
		break
	case 0x9000:
		c.exec9XY0(x, y)
	case 0xA000:
		c.execANNN(addr)
	case 0xB000:
		c.execBNNN(addr)
	case 0xC000:
		c.execCXNN(x, nn)
	case 0xD000:
		c.execDXYN(x, y, n)
	case 0xE000:
		switch c.opcode & 0x00FF {
		case 0x009E:
			c.execEX9E(x)
		case 0x00A1:
			c.execEXA1(x)
		default:
			return errors.New(fmt.Sprintf("E000: Unsupported opcode %v", c.opcode))
		}
		break
	case 0xF000:
		switch c.opcode & 0x0FF {
		case 0x0007:
			c.execFX07(x)
		case 0x000A:
			c.execFX0A(x)
		case 0x0015:
			c.execFX15(x)
		case 0x0018:
			c.execFX18(x)
		case 0x001E:
			c.execFX1E(x)
		case 0x0029:
			c.execFX29(x)
		case 0x0033:
			c.execFX33(x)
		case 0x0055:
			c.execFX55(x)
		case 0x0065:
			c.execFX65(x)
		default:
			return errors.New(fmt.Sprintf("FX00: Unsupported opcode %v", c.opcode))
		}
	default:
		return errors.New(fmt.Sprintf("Unknown opcode: Unsupported opcode %v", c.opcode))
	}

	if c.delayTimer > 0 {
		c.delayTimer--
	}

	if c.SoundTimer > 0 {
		c.SoundTimer--
	}

	return nil
}

func (c *CHIP8) debugDraw() {
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if c.Gfx[x][y] == 1 {
				fmt.Print("0")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Println()
}

func (c *CHIP8) exec0NNN(addr uint16) {
	c.Log("0NNN")
}

func (c *CHIP8) exec00E0() {
	c.Log("00E0")
	c.Gfx = [Width][Height]byte{}
	c.pc += 2
}

func (c *CHIP8) exec00EE() {
	c.Log("00EE")
	c.sp--
	c.pc = c.stack[c.sp]
}

func (c *CHIP8) exec1NNN(addr uint16) {
	c.Log("1NNN")
	c.pc = addr
}

func (c *CHIP8) exec2NNN(addr uint16) {
	c.Log("2NNN")
	c.stack[c.sp] = c.pc + 2
	c.sp++
	c.pc = addr
}

func (c *CHIP8) exec3XNN(x, nn byte) {
	c.Log("3XNN")
	if c.v[x] == nn {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec4XNN(x, nn byte) {
	c.Log("4XNN")
	if c.v[x] != nn {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec5XY0(x, y byte) {
	c.Log("5XY0")
	if c.v[x] == c.v[y] {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec6XNN(x, nn byte) {
	c.Log("6XNN")
	c.v[x] = nn
	c.pc += 2
}

func (c *CHIP8) exec7XNN(x, nn byte) {
	c.Log("7XNN")
	c.v[x] += nn
	c.pc += 2
}

func (c *CHIP8) exec8XY0(x, y byte) {
	c.Log("8XY0")
	c.v[x] = c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY1(x, y byte) {
	c.Log("8XY1")
	c.v[x] |= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY2(x, y byte) {
	c.Log("8XY2")
	c.v[x] &= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY3(x, y byte) {
	c.Log("8XY3")
	c.v[x] ^= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY4(x, y byte) {
	c.Log("8XY4")
	if (c.v[y] + c.v[x]) > 0xFF {
		c.v[0xF] = 0
	} else {
		c.v[0xF] = 1
	}
	c.v[x] += c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY5(x, y byte) {
	c.Log("8XY5")
	if c.v[y] > c.v[x] {
		c.v[0xF] = 0
	} else {
		c.v[0xF] = 1
	}
	c.v[x] -= c.v[y]
	c.pc += 2
}

func (c *CHIP8) exec8XY6(x, y byte) {
	c.Log("8XY6")
	c.v[0xF] = c.v[x] & 0x1
	c.v[x] >>= 1
	c.pc += 2
}

func (c *CHIP8) exec8XY7(x, y byte) {
	c.Log("8XY7")
	if c.v[y] > c.v[x] {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
	c.v[x] = c.v[y] - c.v[x]
	c.pc += 2
}

func (c *CHIP8) exec8XYE(x, y byte) {
	c.Log("8XYE")
	c.v[0xF] = (c.v[x] >> 7) & 0x1
	c.v[x] <<= 1
	c.pc += 2
}

func (c *CHIP8) exec9XY0(x, y byte) {
	c.Log("9XY0")
	if c.v[x] != c.v[y] {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) execANNN(addr uint16) {
	c.Log("ANNN")
	c.i = addr
	c.pc += 2
}

func (c *CHIP8) execBNNN(addr uint16) {
	c.Log("BNNN")
	c.pc = addr + uint16(c.v[0x0])
}

func (c *CHIP8) execCXNN(x, nn byte) {
	c.Log("CXNN")
	b := randomByte()
	c.v[x] = b & nn
	c.pc += 2
}

func (c *CHIP8) execDXYN(x, y, n byte) {
	c.Log("DXYN")
	vx := c.v[x]
	vy := c.v[y]
	var pixel byte

	c.v[0xF] = 0
	for yl := byte(0); yl < n; yl++ { // n = height
		pixel = c.memory[c.i+uint16(yl)]
		for xl := byte(0); xl < 8; xl++ { // width => always 8 pixels

			if (pixel & (0x80 >> xl)) != 0 {
				if c.Gfx[(vx+xl)%64][(vy+yl)%32] == 1 {
					c.v[0xF] = 1
				}
				c.Gfx[(vx+xl)%64][(vy+yl)%32] ^= 1
			}
		}
	}
	c.DrawFlag = true
	c.pc += 2
}

func (c *CHIP8) execEX9E(x byte) {
	c.Log("EX9E")
	if c.keys[c.v[x]] == 1 {
		c.pc += 2
	}
	c.keys[c.v[x]] = 0
	c.pc += 2
}

func (c *CHIP8) execEXA1(x byte) {
	c.Log("EXA1")
	if c.keys[c.v[x]] == 0 {
		c.pc += 2
	}
	c.keys[c.v[x]] = 0
	c.pc += 2
}

func (c *CHIP8) execFX07(x byte) {
	c.Log("FX07")
	c.v[x] = c.delayTimer
	c.pc += 2
}

func (c *CHIP8) execFX0A(x byte) {
	c.Log("FX0A")
	for i, k := range c.keys {
		if k == 1 {
			c.v[x] = byte(i) // value of the key
			c.pc += 2
			return
		}
	}
}

func (c *CHIP8) execFX15(x byte) {
	c.Log("FX15")
	c.delayTimer = c.v[x]
	c.pc += 2
}

func (c *CHIP8) execFX18(x byte) {
	c.Log("FX18")
	c.SoundTimer = c.v[x]
	c.pc += 2
}

func (c *CHIP8) execFX1E(x byte) {
	c.Log("FX1E")
	c.i += uint16(c.v[x])
	c.pc += 2
}

func (c *CHIP8) execFX29(x byte) {
	c.Log("FX29")
	c.i = uint16(5 * c.v[x])
	c.pc += 2
}

func (c *CHIP8) execFX33(x byte) {
	c.Log("FX33")
	c.memory[c.i] = c.v[x] / 100
	c.memory[c.i+1] = (c.v[x] / 10) % 10
	c.memory[c.i+2] = (c.v[x] % 100) % 10
	c.pc += 2
}

func (c *CHIP8) execFX55(x byte) {
	c.Log("FX55")
	for i := byte(0); i <= x; i++ {
		c.memory[c.i+uint16(i)] = c.v[i]
	}
	c.pc += 2
}

func (c *CHIP8) execFX65(x byte) {
	c.Log("FX65")
	for i := byte(0); i <= x; i++ {
		c.v[i] = c.memory[c.i+uint16(i)]
	}
	c.pc += 2
}

func (c *CHIP8) LoadRomHTTP(loc string) error {
	resp, err := http.Get(loc)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)

	for i := 0; i < len(b); i++ {
		c.memory[i+0x200] = b[i]
	}

	return nil
}

func (c *CHIP8) Load(romName string) error {
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
	}
	return nil
}

func (c *CHIP8) Log(opcode string) {
	if c.Logging {
		log.Printf("PC: %v, Executing: %v", c.pc, opcode)
	}
}

func (c *CHIP8) String() string {
	return fmt.Sprintf("PC: %v\n", c.pc)
}

func randomByte() byte {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]
}
