package vm

import (
	"strconv"
)

type OpCode uint32

// opcodes
const (
	OP_LEA OpCode = iota
	OP_IMM
	OP_JMP
	OP_JSR
	OP_BZ
	OP_BNZ
	OP_ENT
	OP_ADJ
	OP_LEV
	OP_LI
	OP_LC
	OP_SI
	OP_SC
	OP_PCH
	OP_OR
	OP_XOR
	OP_AND
	OP_EQ
	OP_NE
	OP_LT
	OP_GT
	OP_GE
	OP_SHL
	OP_SHR
	OP_ADD
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD

	OP_OPEN
	OP_READ
	OP_CLOS
	OP_PRTF
	OP_MALC
	OP_MEST
	OP_MCMP
	OP_MCPY
	OP_MMAP
	OP_DSYM
	OP_QSRT
	OP_EXIT
)

// types
const (
	CHAR Token = iota
	INT
	PTR
)

var opCodeNames = [...]string{
	OP_LEA: "LEA",
	OP_IMM: "IMM",
	OP_JMP: "JMP",
	OP_JSR: "JSR",
	OP_BZ:  "BZ",
	OP_BNZ: "BNZ",
	OP_ENT: "ENT",
	OP_ADJ: "ADJ",
	OP_LEV: "LEV",
	OP_LI:  "LI",
	OP_LC:  "LC",
	OP_SI:  "SI",
	OP_SC:  "SC",
	OP_PCH: "PCH",
	OP_OR:  "OR",
	OP_XOR: "XOR",
	OP_AND: "AND",
	OP_EQ:  "EQ",
	OP_NE:  "NE",
	OP_LT:  "LT",
	OP_GT:  "GT",
	OP_GE:  "GE",
	OP_SHL: "SHL",
	OP_SHR: "SHR",
	OP_ADD: "ADD",
	OP_SUB: "SUB",
	OP_MUL: "MUL",
	OP_DIV: "DIV",
	OP_MOD: "MOD",

	OP_OPEN: "OPEN",
	OP_READ: "READ",
	OP_CLOS: "CLOS",
	OP_PRTF: "PRTF",
	OP_MALC: "MALC",
	OP_MEST: "MEST",
	OP_MCMP: "MCMP",
	OP_MCPY: "MCPY",
	OP_MMAP: "MMAP",
	OP_DSYM: "DSYM",
	OP_QSRT: "QSRT",
	OP_EXIT: "EXIT",
}

func (op OpCode) String() string {
	s := ""
	if op < OpCode(len(opCodeNames)) {
		s = opCodeNames[op]
	}
	if s == "" {
		s = "OpCode(" + strconv.Itoa(int(op)) + ")"
	}
	return s
}
