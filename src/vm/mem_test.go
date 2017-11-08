package vm

import (
	//"fmt"
	"testing"
	"unsafe"
)

func TestOpCode(t *testing.T) {
	testdata := []struct {
		memSize uint64
		index   uint64
		val     OpCode
		ok      bool
		newPos  uint64
	}{
		{100, 0, OP_ADJ, true, uint64(unsafe.Sizeof(OpCode(0)))},
		{100, 1, OP_ENT, true, 1 + uint64(unsafe.Sizeof(OpCode(0)))},
		{100, 50, OP_ENT, true, 50 + uint64(unsafe.Sizeof(OpCode(0)))},
		{100, 100 - uint64(unsafe.Sizeof(OpCode(0))), OP_ENT, true, 100},

		{100, 100 - uint64(unsafe.Sizeof(OpCode(0))) + 1, OP_ENT, false, 0},
		{100, 100, OP_ADJ, false, 0},
	}

	prefix := "TestSetOpCode"

	for i, v := range testdata {

		memView := NewMemView(v.memSize)

		ok, newPos := memView.SetOpCode(v.index, v.val)
		if ok != v.ok {
			t.Errorf("%s[%d] failed: ok = %v, wanted = %v\n", prefix, i, ok, v.ok)
			continue
		}

		if !v.ok {
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: SetOpCode newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}

		val, newPos := memView.GetOpCode(v.index)
		if val != v.val {
			t.Errorf("%s[%d] failed: GetOpCode = %v, wanted = %v\n", prefix, i, val, v.val)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: GetOpCode newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}
	}
}

func TestUint64(t *testing.T) {
	testdata := []struct {
		memSize uint64
		index   uint64
		val     uint64
		ok      bool
		newPos  uint64
	}{
		{100, 0, 789, true, uint64(unsafe.Sizeof(uint64(0)))},
		{100, 1, 121349, true, 1 + uint64(unsafe.Sizeof(uint64(0)))},
		{100, 50, 12876433, true, 50 + uint64(unsafe.Sizeof(uint64(0)))},
		{100, 100 - uint64(unsafe.Sizeof(uint64(0))), 9999888, true, 100},

		{100, 100 - uint64(unsafe.Sizeof(uint64(0))) + 1, 56234496, false, 0},
		{100, 100, 1235678, false, 0},
	}

	prefix := "TestUint64"

	for i, v := range testdata {

		memView := NewMemView(v.memSize)

		ok, newPos := memView.SetUint64(v.index, v.val)
		if ok != v.ok {
			t.Errorf("%s[%d] failed: ok = %v, wanted = %v\n", prefix, i, ok, v.ok)
			continue
		}

		if !v.ok {
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: SetUint64 newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}

		val, newPos := memView.GetUint64(v.index)
		if val != v.val {
			t.Errorf("%s[%d] failed: GetUint64 = %v, wanted = %v\n", prefix, i, val, v.val)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: GetUint64 newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}
	}
}

func TestInt64(t *testing.T) {
	testdata := []struct {
		memSize uint64
		index   uint64
		val     int64
		ok      bool
		newPos  uint64
	}{
		{100, 0, 789, true, uint64(unsafe.Sizeof(int64(0)))},
		{100, 1, -121349, true, 1 + uint64(unsafe.Sizeof(int64(0)))},
		{100, 50, 12876433, true, 50 + uint64(unsafe.Sizeof(int64(0)))},
		{100, 100 - uint64(unsafe.Sizeof(int64(0))), 9999888, true, 100},

		{100, 100 - uint64(unsafe.Sizeof(int64(0))) + 1, 56234496, false, 0},
		{100, 100, -1235678, false, 0},
	}

	prefix := "TestInt64"

	for i, v := range testdata {

		memView := NewMemView(v.memSize)

		ok, newPos := memView.SetInt64(v.index, v.val)
		if ok != v.ok {
			t.Errorf("%s[%d] failed: ok = %v, wanted = %v\n", prefix, i, ok, v.ok)
			continue
		}

		if !v.ok {
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: SetUint64 newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}

		val, newPos := memView.GetInt64(v.index)
		if val != v.val {
			t.Errorf("%s[%d] failed: GetInt64 = %v, wanted = %v\n", prefix, i, val, v.val)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: GetInt64 newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}
	}
}

func TestChar(t *testing.T) {
	testdata := []struct {
		memSize uint64
		index   uint64
		val     byte
		ok      bool
		newPos  uint64
	}{
		{100, 0, 'a', true, uint64(unsafe.Sizeof(byte(0)))},
		{100, 1, 'A', true, 1 + uint64(unsafe.Sizeof(byte(0)))},
		{100, 50, '0', true, 50 + uint64(unsafe.Sizeof(byte(0)))},
		{100, 100 - uint64(unsafe.Sizeof(byte(0))), '9', true, 100},

		{100, 100 - uint64(unsafe.Sizeof(byte(0))) + 1, ';', false, 0},
		{100, 100, '[', false, 0},
	}

	prefix := "TestChar"

	for i, v := range testdata {

		memView := NewMemView(v.memSize)

		ok, newPos := memView.SetChar(v.index, v.val)
		if ok != v.ok {
			t.Errorf("%s[%d] failed: ok = %v, wanted = %v\n", prefix, i, ok, v.ok)
			continue
		}

		if !v.ok {
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: SetChar newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}

		val, newPos := memView.GetChar(v.index)
		if val != v.val {
			t.Errorf("%s[%d] failed: GetChar = %v, wanted = %v\n", prefix, i, val, v.val)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: GetChar newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}
	}
}

func TestString(t *testing.T) {
	testdata := []struct {
		memSize uint64
		index   uint64
		val     string
		ok      bool
		newPos  uint64
	}{
		{100, 0, "abc", true, uint64(len("abc")) + 1},
		{100, 1, "tsde", true, 1 + uint64(len("tsde")) + 1},
		{100, 50, "123adafd", true, 50 + uint64(len("123adafd")) + 1},
		{100, 100 - uint64(len("h2138zx")) - 1, "h2138zx", true, 100},

		{100, 100 - uint64(len("h2138zx")), "h2138zx", false, 100},
		{100, 100, "h2138zx", false, 0},
	}

	prefix := "TestString"

	for i, v := range testdata {

		memView := NewMemView(v.memSize)

		ok, newPos := memView.SetString(v.index, v.val)
		if ok != v.ok {
			t.Errorf("%s[%d] failed: ok = %v, wanted = %v\n", prefix, i, ok, v.ok)
			continue
		}

		if !v.ok {
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: SetString newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}

		val, newPos := memView.GetString(v.index)
		if val != v.val {
			t.Errorf("%s[%d] failed: GetString = %v, wanted = %v\n", prefix, i, val, v.val)
			continue
		}

		if newPos != v.newPos {
			t.Errorf("%s[%d] failed: GetString newPos = %v, wanted = %v\n", prefix, i, newPos, v.newPos)
			continue
		}
	}
}
