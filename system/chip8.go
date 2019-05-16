package system

import (
	"bufio"
	"fmt"
	"os"
)

var (
	soundTimer int
	delayTimer int
)

type CHIP8 struct {
	cpu
	gfx      [64 * 32]byte // display
	key      [16]byte      // current key state
	DrawFlag bool
}

type cpu struct {
	opcode uint16 // opcode is 2 bytes
	memory [4096]byte
	v      [16]byte // CPU registers
	i      uint16   // index register
	pc     uint16   // program counter
	stack  [16]uint16
	sp     uint16 // stackpointer
}

func (c *CHIP8) Initialize() {
	c.pc = 0x200 // = 512: initial point where a program will start
	c.opcode = 0x00
	c.i = 0x00
	c.sp = 0x00
	c.DrawFlag = true

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
func (c *CHIP8) Cycle() {
	c.opcode = uint16(c.memory[c.pc])
	c.opcode <<= 8
	c.opcode |= uint16(c.memory[c.pc+1])

	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode & 0x000F {
		case 0x0000:
			c.exec00E0()
			break
		case 0x000E:
			c.exec00EE()
			break
		}
	case 0x1000:
		c.exec1NNN()
		break
	case 0x2000:
		c.exec2NNN()
		break
	case 0x3000:
		c.exec3XNN()
		break
	case 0x4000:
		c.exec4XNN()
		break
	case 0x5000:
		c.exec5XY0()
		break
	case 0x6000:
		c.exec6XNN()
		break
	case 0x7000:
		c.exec7XNN()
		break
	case 0x0004:
		c.exec8XY4()
		break
	case 0xA000:
		c.execANNN()
		break
	case 0x0033:
		c.execFX33()
		break
	default:
		//fmt.Println("Unsupported instruction")
	}

	//fmt.Printf("%s\n", fmt.Sprintf("%x", c.opcode))
	fmt.Printf("%s\n", fmt.Sprintf("%x", c.pc))
}

func (c *CHIP8) exec00E0() {
	fmt.Printf("Executing 00E0\n")
	c.gfx = [64 * 32]byte{}
	c.pc += 2
}

func (c *CHIP8) exec00EE() {
	fmt.Printf("Executing 00EE\n")
	c.pc = c.stack[c.sp]
	c.sp--
}

func (c *CHIP8) exec1NNN() {
	fmt.Printf("Executing 1NNN\n")
	c.pc = c.opcode & 0x0FFF
}

func (c *CHIP8) exec2NNN() {
	fmt.Printf("Executing 2NNN\n")
	c.stack[c.sp] = c.pc
	c.sp++
	c.pc = c.opcode & 0x0FFF
}

func (c *CHIP8) exec3XNN() {
	fmt.Printf("Executing 3XNN\n")
	if uint16(c.v[(c.opcode>>8)&0x0F00]) == c.opcode&0x00FF {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec4XNN() {
	fmt.Printf("Executing 5XY0\n")
	if uint16(c.v[c.opcode&0x0F00]) != c.opcode&0x00FF {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec5XY0() {
	fmt.Printf("Executing 5XY0\n")
	if c.v[(c.opcode>>8)&0x0F00] == c.v[(c.opcode>>4)&0x00F0] {
		c.pc += 4
		return
	}
	c.pc += 2
}

func (c *CHIP8) exec6XNN() {
	fmt.Printf("Executing 6XNN\n")
	c.v[(c.opcode>>8)&0x0F00] = byte(c.opcode & 0x00FF)
	c.pc += 2
}

func (c *CHIP8) exec7XNN() {
	fmt.Println("Executing 7XNN")
	c.v[(c.opcode>>8)&0x0F00] += byte(c.opcode & 0x00FF)
	c.pc += 2
}

func (c *CHIP8) exec8XY4() {
	fmt.Printf("Executing 8XY4\n")
	if c.v[(c.opcode&0x00F0)>>4] > (0xFF - c.v[(c.opcode&0x0F00)>>8]) {
		c.v[0xF] = 1
	} else {
		c.v[0xF] = 0
	}
	c.v[(c.opcode&0x0F00)>>8] += c.v[(c.opcode&0x00F0)>>4]
	c.pc += 2
}

func (c *CHIP8) execANNN() {
	fmt.Printf("Executing ANNN\n")
	c.i = c.opcode & 0x0FFF
	c.pc += 2
}

func (c *CHIP8) execFX33() {
	c.memory[c.i] = c.v[(c.opcode&0x0F00)>>8] / 100
	c.memory[c.i+1] = (c.v[(c.opcode&0x0F00)>>8] / 10) % 10
	c.memory[c.i+2] = (c.v[(c.opcode&0x0F00)>>8] % 100) % 10
	c.pc += 2
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
		//hex := fmt.Sprintf("%x", c.memory[i+0x200])
		//fmt.Printf("%s\n", hex)
	}

	return nil
}
