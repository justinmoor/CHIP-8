package system

import (
	"bufio"
	"fmt"
	"os"
)

var (
	soundTimer 	int
	delayTimer	int
)

type CHIP8 struct{
	cpu
	gfx			[64 * 32]byte	// display
	key			[16]byte		// current key state
	DrawFlag	bool
}

type cpu struct {
	opcode 	uint16		// opcode is 2 bytes
	memory 	[4096]byte
	v		[16]byte	// CPU registers
	i		uint16		// index register
	pc		uint16		// program counter
	stack	[16]uint16
	sp		uint16		// stackpointer
}

func (c *CHIP8) Initialize(){
	c.pc		= 0x200 // = 512: initial point where a program will start
	c.opcode	= 0x00
	c.i			= 0x00
	c.sp		= 0x00
	c.DrawFlag	= true

	// clear all the memory
	c.memory	= [4096]byte{}
	c.v			= [16]byte{}
	c.stack		= [16]uint16{}

	// load font into memory
	for i := 0; i < 0x50; i++ {
		c.memory[i] = font[i]
	}
}

// fetch opcode
// decode opcode
// execude opcode
// update timers
func (c *CHIP8) Cycle(){
	//c.opcode = uint16((c.memory[c.pc] << 8) | c.memory[c.pc + 1])
	// compiler optimization:
	c.opcode = uint16(c.memory[c.pc])
	c.opcode <<= 8
	c.opcode |= uint16(c.memory[c.pc + 1])
	c.pc += 2
	/*
	int a = opcode & 0xF000;
	int b = opcode & 0x0FFF;
	if(a == 0x1000)
	   doJump(b);
	*/
	switch c.opcode & 0xF000 {
	case 0x0000:
		switch c.opcode & 0x000F {
		case 0x0000:
			c.clearDisplay()
			break
		case 0x000E:
			c.pc = c.stack[c.sp]
			c.sp--
			break
		}
	case 0x1000:
		c.pc = c.opcode & 0x0FFF
		break
	case 0x2000:
		c.stack[c.sp] = c.pc
		c.sp++
		c.pc = c.opcode & 0x0FFF
		break
	case 0x3000:
		if uint16(c.v[c.opcode & 0x0F00]) == c.opcode & 0x00FF {
			c.pc += 2
		}
		break
	case 0x0004:
		if c.v[(c.opcode & 0x00F0) >> 4] > (0xFF - c.v[(c.opcode & 0x0F00) >> 8]) {
			c.v[0xF] = 1
		} else {
			c.v[0xF] = 0
		}
		c.v[(c.opcode & 0x0F00) >> 8] += c.v[(c.opcode & 0x00F0) >> 4]
		c.pc += 2
		break
	case 0xA000:
		c.i = c.opcode & 0x0FFF
		c.pc += 2
		break
	case 0x0033:
		c.memory[c.i]		= c.v[(c.opcode & 0x0F00) >> 8] / 100
		c.memory[c.i + 1]	= (c.v[(c.opcode & 0x0F00) >> 8] / 10) % 10
		c.memory[c.i + 2]	= (c.v[(c.opcode & 0x0F00) >> 8] % 100) % 10
		break
	default:
		fmt.Println("Unsupported")
	}
}

func (c *CHIP8) clearDisplay(){
	c.gfx = [64 * 32]byte{}
}

func (c *CHIP8) Load(romName string) error{
	rom, err := os.Open(romName)
	if err != nil {
		return err
	}
	defer rom.Close()

	stats, err := rom.Stat()
	if err != nil {
		return err
	}

	bytes := make([]byte, stats.Size())
	buffer := bufio.NewReader(rom)

	if _, err = buffer.Read(bytes); err != nil {
		return err
	}

	for i := 0; i < len(bytes); i++ {
		c.memory[i + 0x200] = bytes[i]
		hex := fmt.Sprintf("%x", c.memory[i+0x200])
		fmt.Printf("%s\n", hex)
	}

	return nil
}
