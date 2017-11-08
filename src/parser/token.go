package parser

import (
	_ "fmt"
)

type Token int32

// tokens and classes (operators last and in precedence order)
const (
	TOKEN_EOF Token = 1
	TOKEN_NUM Token = iota + 128
	TOKEN_FUN
	TOKEN_SYS
	TOKEN_GLO
	TOKEN_LOC
	TOKEN_ID
	TOKEN_CHAR
	TOKEN_STRING
	TOKEN_ELSE
	TOKEN_ENUM
	TOKEN_IF
	TOKEN_INT
	TOKEN_RETURN
	TOKEN_SIZEOF
	TOKEN_WHILE

	TOKEN_ASSIGN
	TOKEN_COND
	TOKEN_LOR
	TOKEN_LAN
	TOKEN_OR
	TOKEN_XOR
	TOKEN_AND
	TOKEN_EQ
	TOKEN_NE
	TOKEN_LT
	TOKEN_GT
	TOKEN_LE
	TOKEN_GE
	TOKEN_SHL
	TOKEN_SHR
	TOKEN_ADD
	TOKEN_SUB
	TOKEN_MUL
	TOKEN_DIV
	TOKEN_MOD
	TOKEN_INC
	TOKEN_DEC
	TOKEN_BRAK
)

var tokenNames = [...]string{
	TOKEN_EOF:    "EOF",
	TOKEN_NUM:    "NUM",
	TOKEN_FUN:    "FUN",
	TOKEN_SYS:    "SYS",
	TOKEN_GLO:    "GLO",
	TOKEN_LOC:    "LOC",
	TOKEN_ID:     "ID",
	TOKEN_CHAR:   "CHAR",
	TOKEN_STRING: "STRING",
	TOKEN_ELSE:   "ELSE",
	TOKEN_ENUM:   "ENUM",
	TOKEN_IF:     "IF",
	TOKEN_INT:    "INT",
	TOKEN_RETURN: "RETURN",
	TOKEN_SIZEOF: "SIZEOF",
	TOKEN_WHILE:  "WHILE",

	TOKEN_ASSIGN: "ASSIGN",
	TOKEN_COND:   "COND",
	TOKEN_LOR:    "LOR",
	TOKEN_LAN:    "LAN",
	TOKEN_OR:     "OR",
	TOKEN_XOR:    "XOR",
	TOKEN_AND:    "AND",
	TOKEN_EQ:     "EQ",
	TOKEN_NE:     "NE",
	TOKEN_LT:     "LT",
	TOKEN_GT:     "GT",
	TOKEN_LE:     "LE",
	TOKEN_GE:     "GE",
	TOKEN_SHL:    "SHL",
	TOKEN_SHR:    "SHR",
	TOKEN_ADD:    "ADD",
	TOKEN_SUB:    "SUB",
	TOKEN_MUL:    "MUL",
	TOKEN_DIV:    "DIV",
	TOKEN_MOD:    "MOD",
	TOKEN_INC:    "INC",
	TOKEN_DEC:    "DEC",
	TOKEN_BRAK:   "BRAK",
}

func (token Token) String() string {
	s := ""
	if token < Token(len(tokenNames)) {
		s = tokenNames[token]
	}
	if s == "" {
		s = PrintChar(byte(token))
	}
	return s
}
