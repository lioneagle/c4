package vm

import (
	"config"
	"fmt"
	//"reflect"
	"unsafe"
)

type Registers struct {
	pc    uint64
	sp    uint64
	bp    uint64
	ax    uint64
	cycle uint64
}

type Vim struct {
	text      *memView
	data      *memView
	stack     *memView
	reg       Registers
	textSize  uint64
	dataSize  uint64
	stackSize uint64
}

const (
	VM_DEFAULT_TEXT_SIZE  = 256 * 1024
	VM_DEFAULT_DATA_SIZE  = 256 * 1024
	VM_DEFAULT_STACK_SIZE = 256 * 1024
)

func SizeofUint64() uint64 {
	return uint64(unsafe.Sizeof(uint64(0)))
}

func SizeofInt64() uint64 {
	return uint64(unsafe.Sizeof(int64(0)))
}

func SizeofOpCode() uint64 {
	return uint64(unsafe.Sizeof(OpCode(0)))
}

func SizeofChar() uint64 {
	return uint64(unsafe.Sizeof(byte(0)))
}

func NewVM(textSize, dataSize, stackSize uint64) *Vim {
	vim := &Vim{}
	vim.text = NewMemView(textSize + dataSize + stackSize)
	vim.data = &memView{data: vim.text.data[textSize : textSize+dataSize]}
	vim.stack = &memView{data: vim.text.data[textSize+dataSize:]}

	vim.reg.sp = uint64(len(vim.text.data))
	vim.reg.bp = vim.reg.sp

	vim.textSize = textSize
	vim.dataSize = dataSize
	vim.stackSize = stackSize

	return vim
}

func (this *Vim) DataBegin() uint64 {
	return this.textSize
}

func (this *Vim) AddOpCode(pos uint64, op OpCode) (ok bool, newPos uint64) {
	return this.text.SetOpCode(pos, op)
}

func (this *Vim) AddUint64(pos uint64, val uint64) (ok bool, newPos uint64) {
	return this.text.SetUint64(pos, val)
}

func (this *Vim) AddChar(pos uint64, val byte) (ok bool, newPos uint64) {
	return this.text.SetChar(pos, val)
}

func (this *Vim) AddDataString(pos uint64, str string) (ok bool, newPos uint64) {
	return this.text.SetString(pos, str)
}

func (this *Vim) OpCode(pos uint64) OpCode {
	ret, _ := this.text.GetOpCode(pos)
	return ret
}

func (this *Vim) Uint64(pos uint64) uint64 {
	ret, _ := this.text.GetUint64(pos)
	return ret
}

func (this *Vim) CurrentOpCode() OpCode {
	return this.OpCode(this.reg.pc)
}

func (this *Vim) PopOpCode() (val OpCode) {
	val, this.reg.pc = this.text.GetOpCode(this.reg.pc)
	return val
}

func (this *Vim) incPC() {
	this.reg.pc += SizeofUint64()
}

func (this *Vim) pcAddrValue() (val uint64) {
	val, _ = this.text.GetUint64(this.reg.pc)
	return val
}

func (this *Vim) popPCAddrValue() (val uint64) {
	val, this.reg.pc = this.text.GetUint64(this.reg.pc)
	return val
}

func (this *Vim) popStack() (val uint64) {
	val, this.reg.sp = this.text.GetUint64(this.reg.sp)
	return val
}

func (this *Vim) topStack() (val uint64) {
	val, _ = this.text.GetUint64(this.reg.sp)
	return val
}

func (this *Vim) PushStack(val uint64) (ok bool) {
	this.reg.sp -= SizeofUint64()
	ok, _ = this.text.SetUint64(this.reg.sp, val)
	return ok
}

func (this *Vim) Run(config *config.RunConfig) uint64 {
	this.reg.cycle = 0

	for {
		instruct := this.PopOpCode()
		this.reg.cycle++
		if config.Debug {
			fmt.Printf("%d> %.4s", this.reg.cycle, instruct.String())
			if instruct < OP_ADJ {
				fmt.Printf(" %d", this.CurrentOpCode())
			}
			fmt.Println()
		}

		switch instruct {
		case OP_IMM: // load immediate value to ax
			this.reg.ax = this.popPCAddrValue()

		case OP_LC: // load character to ax, address in ax
			val, _ := this.text.GetChar(this.reg.ax)
			this.reg.ax = uint64(val)

		case OP_LI: // load integer to ax, address in ax
			this.reg.ax, _ = this.text.GetUint64(this.reg.ax)

		case OP_SC: // save character to address, value in ax, address on stack
			addr := this.popStack()
			this.text.SetChar(addr, byte(this.reg.ax))

		case OP_SI: // save integer to address, value in ax, address on stack
			addr := this.popStack()
			this.text.SetUint64(addr, this.reg.ax)

		case OP_PUSH: // push the value of ax onto the stack
			this.PushStack(this.reg.ax)

		case OP_JMP: // jump to the address
			this.reg.pc = this.pcAddrValue()

		case OP_JZ: // jump if ax is zero
			if this.reg.ax == 0 {
				this.reg.pc = this.pcAddrValue()
			} else {
				this.incPC()
			}

		case OP_JNZ: // jump if ax is not zero
			if this.reg.ax != 0 {
				this.reg.pc = this.pcAddrValue()
			} else {
				this.incPC()
			}

		case OP_CALL: // call subroutine
			addr := this.reg.pc
			this.incPC()
			this.PushStack(this.reg.pc)
			this.reg.pc = addr
			this.reg.pc = this.pcAddrValue()

		case OP_ENT: // enter subroutine,  make new stack frame
			this.PushStack(this.reg.bp)
			this.reg.bp = this.reg.sp
			this.reg.sp -= this.popPCAddrValue()

		case OP_ADJ: // stack adjust, add esp, <size>
			this.reg.sp += this.popPCAddrValue()

		case OP_LEV: // leave subroutine, restore call frame and PC
			this.reg.sp = this.reg.bp
			this.reg.bp = this.popStack()
			this.reg.pc = this.popStack()

		case OP_LEA: // load address for arguments
			this.reg.ax = this.reg.bp + this.popPCAddrValue()

		case OP_OR:
			this.reg.ax |= this.popStack()

		case OP_XOR:
			this.reg.ax ^= this.popStack()

		case OP_AND:
			this.reg.ax &= this.popStack()

		case OP_EQ:
			if this.popStack() == this.reg.ax {
				this.reg.ax = 1
			} else {
				this.reg.ax = 0
			}

		case OP_NE:
			if this.popStack() != this.reg.ax {
				this.reg.ax = 1
			} else {
				this.reg.ax = 0
			}

		case OP_LT:
			if this.popStack() < this.reg.ax {
				this.reg.ax = 1
			} else {
				this.reg.ax = 0
			}

		case OP_LE:
			if this.popStack() <= this.reg.ax {
				this.reg.ax = 1
			} else {
				this.reg.ax = 0
			}

		case OP_GT:
			if this.popStack() > this.reg.ax {
				this.reg.ax = 1
			} else {
				this.reg.ax = 0
			}

		case OP_GE:
			if this.popStack() >= this.reg.ax {
				this.reg.ax = 1
			} else {
				this.reg.ax = 0
			}

		case OP_SHL:
			this.reg.ax = this.popStack() << this.reg.ax

		case OP_SHR:
			this.reg.ax = this.popStack() >> this.reg.ax

		case OP_ADD:
			this.reg.ax = this.popStack() + this.reg.ax

		case OP_SUB:
			this.reg.ax = this.popStack() - this.reg.ax

		case OP_MUL:
			this.reg.ax = this.popStack() * this.reg.ax

		case OP_DIV:
			this.reg.ax = this.popStack() / this.reg.ax

		case OP_MOD:
			this.reg.ax = this.popStack() % this.reg.ax

		case OP_EXIT:
			fmt.Printf("exit(%d)\n", this.topStack())
			return this.topStack()
		}
	}
}
