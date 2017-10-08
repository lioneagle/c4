package parser

type Token int32

// tokens and classes (operators last and in precedence order)
const (
	NUM Token = iota + 128
	FUN
	SYS
	GLO
	LOC
	ID
	CHAR
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

func (this Token) String() string {
	return ""
}
