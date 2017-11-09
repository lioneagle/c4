package vm

import (
	//"fmt"
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
	return *(*OpCode)(unsafe.Pointer(&this.data[pos])), pos + SizeofOpCode()

}

func (this *memView) SetOpCode(pos uint64, val OpCode) (ok bool, newPos uint64) {
	if (pos + SizeofOpCode()) > uint64(len(this.data)) {
		return false, pos
	}
	*(*OpCode)(unsafe.Pointer(&this.data[pos])) = val
	return true, pos + SizeofOpCode()
}

func (this *memView) GetChar(pos uint64) (val byte, newPos uint64) {
	return this.data[pos], pos + SizeofChar()
}

func (this *memView) SetChar(pos uint64, val byte) (ok bool, newPos uint64) {
	if pos >= uint64(len(this.data)) {
		return false, pos
	}
	this.data[pos] = val
	return true, pos + 1
}

func (this *memView) GetInt64(pos uint64) (val int64, newPos uint64) {
	return *(*int64)(unsafe.Pointer(&this.data[pos])), pos + SizeofInt64()
}

func (this *memView) SetInt64(pos uint64, val int64) (ok bool, newPos uint64) {
	if (pos + SizeofInt64()) > uint64(len(this.data)) {
		return false, pos
	}
	*(*int64)(unsafe.Pointer(&this.data[pos])) = val
	return true, pos + SizeofInt64()
}

func (this *memView) GetUint64(pos uint64) (val uint64, newPos uint64) {
	return *(*uint64)(unsafe.Pointer(&this.data[pos])), pos + SizeofUint64()
}

func (this *memView) SetUint64(pos uint64, val uint64) (ok bool, newPos uint64) {
	if (pos + SizeofUint64()) > uint64(len(this.data)) {
		return false, pos
	}
	*(*uint64)(unsafe.Pointer(&this.data[pos])) = val
	return true, pos + SizeofUint64()
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
	if pos+uint64(len(val)) > uint64(len(this.data)) {
		return false, pos
	}
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
