package vm

import (
	"config"
	"fmt"
	//"reflect"
	"unsafe"
)

type memView struct {
	data []byte
}

func NewMemView(size uint64) *memView {
	return &memView{data: make([]byte, size)}
}

func (this *memView) GetOpCode(pos uint64) (val OpCode, newPos uint64) {
	return *(*OpCode)(unsafe.Pointer(&this.data[pos])), pos + uint64(unsafe.Sizeof(val))

}

func (this *memView) SetOpCode(pos uint64, val OpCode) (ok bool, newPos uint64) {
	if (pos + uint64(unsafe.Sizeof(val))) > uint64(len(this.data)) {
		return false, pos
	}
	*(*OpCode)(unsafe.Pointer(&this.data[pos])) = val
	return true, pos + uint64(unsafe.Sizeof(val))
}

func (this *memView) GetChar(pos uint64) (val byte, newPos uint64) {
	return this.data[pos], pos + uint64(unsafe.Sizeof(val))
}

func (this *memView) SetChar(pos uint64, val byte) (ok bool, newPos uint64) {
	if pos >= uint64(len(this.data)) {
		return false, pos
	}
	this.data[pos] = val
	return true, pos + 1
}

func (this *memView) GetInt64(pos uint64) (val int64, newPos uint64) {
	return *(*int64)(unsafe.Pointer(&this.data[pos])), pos + uint64(unsafe.Sizeof(val))
}

func (this *memView) SetInt64(pos uint64, val int64) (ok bool, newPos uint64) {
	if (pos + uint64(unsafe.Sizeof(val))) > uint64(len(this.data)) {
		return false, pos
	}
	*(*int64)(unsafe.Pointer(&this.data[pos])) = val
	return true, pos + uint64(unsafe.Sizeof(val))
}

func (this *memView) GetUint64(pos uint64) (val uint64, newPos uint64) {
	return *(*uint64)(unsafe.Pointer(&this.data[pos])), pos + uint64(unsafe.Sizeof(val))
}

func (this *memView) SetUint64(pos uint64, val uint64) (ok bool, newPos uint64) {
	if (pos + uint64(unsafe.Sizeof(val))) > uint64(len(this.data)) {
		return false, pos
	}
	*(*uint64)(unsafe.Pointer(&this.data[pos])) = val
	return true, pos + uint64(unsafe.Sizeof(val))
}

func (this *memView) GetString(pos uint64) (val string, newPos uint64) {
	var str []byte
	p := uintptr(unsafe.Pointer(&this.data[pos]))
	end := uintptr(unsafe.Pointer(&this.data[0])) + uintptr(len(this.data))
	for p < end {
		ch := *((*byte)(unsafe.Pointer(p)))
		if ch == 0 {
			break
		}
		str = append(str, ch)
		p++
	}
	return string(str), pos + uint64(len(str)) + 1
}

func (this *memView) SetString(pos uint64, val string) (ok bool, newPos uint64) {
	p := uintptr(unsafe.Pointer(&this.data[pos]))
	end := uintptr(unsafe.Pointer(&this.data[0])) + uintptr(len(this.data))
	i := 0
	for p < end && i < len(val) {
		*((*byte)(unsafe.Pointer(p))) = val[i]
		p++
		i++
	}
	if p >= end {
		return false, pos
	}
	*((*byte)(unsafe.Pointer(p))) = 0
	return true, pos + uint64(i) + 1
}

type Registers struct {
	pc    uint64
	sp    uint64
	bp    uint64
	ax    uint64
	cycle uint64
}

type Vim struct {
	text  *memView
	data  *memView
	stack *memView
	reg   Registers
}

const (
	VM_TEXT_SIZE  = 256 * 1024
	VM_DATA_SIZE  = 256 * 1024
	VM_STACK_SIZE = 256 * 1024
)

func NewVM(textSize, dataSize, stackSize uint64) *Vim {
	vim := &Vim{}
	vim.text = NewMemView(textSize + dataSize + stackSize)
	vim.data = &memView{data: vim.text.data[textSize : textSize+dataSize]}
	vim.stack = &memView{data: vim.text.data[textSize+dataSize:]}

	vim.reg.sp = uint64(len(vim.text.data))
	vim.reg.bp = vim.reg.sp

	return vim
}

func (this *Vim) addOpCode(pos uint64, op OpCode) (ok bool, newPos uint64) {
	return this.text.SetOpCode(pos, op)
}

func (this *Vim) addUint64(pos uint64, val uint64) (ok bool, newPos uint64) {
	return this.text.SetUint64(pos, val)
}

func (this *Vim) addChar(pos uint64, val byte) (ok bool, newPos uint64) {
	return this.text.SetChar(pos, val)
}

func (this *Vim) addDataString(pos uint64, str string) (ok bool, newPos uint64) {
	return this.text.SetString(pos, str)
}

func (this *Vim) OpCode(pos uint64) OpCode {
	ret, _ := this.text.GetOpCode(pos)
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
	this.reg.pc += uint64(unsafe.Sizeof(uint64(0)))
}

func (this *Vim) pcAddrValue() (val uint64) {
	val, _ = this.text.GetUint64(this.reg.pc)
	return val
}

func (this *Vim) popPCAddrValue() (val uint64) {
	val, this.reg.pc = this.text.GetUint64(this.reg.pc)
	return val
}

func (this *Vim) popStack() uint64 {
	this.reg.sp -= uint64(unsafe.Sizeof(uint64(0)))
	val, _ := this.stack.GetUint64(this.reg.sp)
	return val
}

func (this *Vim) topStack() (val uint64) {
	val, _ = this.stack.GetUint64(this.reg.sp - uint64(unsafe.Sizeof(val)))
	return val
}

func (this *Vim) PushStack(val uint64) (ok bool) {
	ok, this.reg.sp = this.text.SetUint64(this.reg.sp, val)
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
