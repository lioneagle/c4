package parser

import (
	"config"
)

type Symbol struct {
	Token  Token
	Name   string
	Class  int32
	Type   int32
	HClass int32
	HType  int32
	HValue int32
	Id     string
}

type Symbols []*Symbol

type Parser struct {
	src             []byte
	pos             int
	localPos        int
	emitPos         int
	localEmitPos    int
	line            int
	currentExprType int32
	token           Token
	tokenValue      int32
	symbols         Symbols
}

func (this *Parser) Next(config *config.RunConfig) {
	for this.pos < len(this.src) {
		ch := this.src[this.pos]

		if ch == '\n' {
			if config.PrintSource {
				this.localPos = this.pos
			}
			this.line++
		}
	}
}
