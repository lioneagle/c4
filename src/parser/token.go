package parser

import (
	_ "fmt"
)

type Token int32

// tokens and classes (operators last and in precedence order)
const (
	EOF Token = 1
	NUM Token = iota + 128
	FUN
	SYS
	GLO
	LOC
	ID
	CHAR
	STRING
	ELSE
	ENUM
	IF
	INT
	RETURN
	SIZEOF
	WHILE

	ASSIGN
	COND
	LOR
	LAN
	OR
	XOR
	AND
	EQ
	NE
	LT
	GT
	LE
	GE
	SHL
	SHR
	ADD
	SUB
	MUL
	DIV
	MOD
	INC
	DEC
	BRAK
)

var tokenNames = [...]string{
	EOF:    "EOF",
	NUM:    "NUM",
	FUN:    "FUN",
	SYS:    "SYS",
	GLO:    "GLO",
	LOC:    "LOC",
	ID:     "ID",
	CHAR:   "CHAR",
	STRING: "STRING",
	ELSE:   "ELSE",
	ENUM:   "ENUM",
	IF:     "IF",
	INT:    "INT",
	RETURN: "RETURN",
	SIZEOF: "SIZEOF",
	WHILE:  "WHILE",

	ASSIGN: "ASSIGN",
	COND:   "COND",
	LOR:    "LOR",
	LAN:    "LAN",
	OR:     "OR",
	XOR:    "XOR",
	AND:    "AND",
	EQ:     "EQ",
	NE:     "NE",
	LT:     "LT",
	GT:     "GT",
	LE:     "LE",
	GE:     "GE",
	SHL:    "SHL",
	SHR:    "SHR",
	ADD:    "ADD",
	SUB:    "SUB",
	MUL:    "MUL",
	DIV:    "DIV",
	MOD:    "MOD",
	INC:    "INC",
	DEC:    "DEC",
	BRAK:   "BRAK",
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
