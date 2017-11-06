package vm

import (
	//"fmt"

	"config"
	"testing"
)

func TestAdd(t *testing.T) {
	config := &config.RunConfig{}
	config.Debug = true
	vim := NewVM(100, 100, 100)

	i := uint64(0)
	_, i = vim.addOpCode(i, OP_IMM)
	_, i = vim.addUint64(i, 10)
	_, i = vim.addOpCode(i, OP_PUSH)
	_, i = vim.addOpCode(i, OP_IMM)
	_, i = vim.addUint64(i, 20)
	_, i = vim.addOpCode(i, OP_ADD)
	_, i = vim.addOpCode(i, OP_PUSH)
	_, i = vim.addOpCode(i, OP_EXIT)

	ret := vim.Run(config)
	if ret != 30 {
		t.Errorf("TestAdd failed, ret = %v, wanted = %v\n", ret, 30)
	}
}
