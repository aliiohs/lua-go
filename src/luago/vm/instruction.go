package vm

const MAXARG_Bx = 1<<18 - 1
const MAXARG_sBx = MAXARG_Bx >> 1

type Instruction uint32

func (self Instruction) Opcode() int {
	// 获取低6位
	return int(self & 0x3F)
}

func (self Instruction) ABC() (a, b, c int) {
	//8
	a = int(self >> 6 & 0xFF)
	//9
	b = int(self >> 14 & 0x1FF)
	//9
	c = int(self >> 23 & 0x1FF)
	return
}

func (self Instruction) ABx() (a, b int) {

	a = int(self >> 6 & 0xFF)
	b = int(self >> 14)
	return
}

func (self Instruction) AsBx() (a, sbx int) {

	a, bx := self.ABx()
	return a, bx - MAXARG_Bx
}

func (self Instruction) Ax() int {

	return int(self >> 6)
}

func (self Instruction) OpName() string {

	return opcodes[self.Opcode()].name

}

func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}
