package parser

import (
	"config"
	"fmt"
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

type Symbols [string]*Symbol

type Parser struct {
	src             []byte
	pos             int
	lastPos         int
	emitPos         int
	lastEmitPos     int
	text            []OpCode
	data            []int
	line            int
	currentExprType int32
	token           Token
	tokenValue      int32
	symbols         Symbols
	ival            int
}

func (this *Parser) Next(config *config.RunConfig) {
	for this.pos < len(this.src) {
		ch := this.src[this.pos]
		this.pos++

		if ch == '\n' {
			if config.PrintSource {
				fmt.Printf("%d: %s", this.line, this.src[this.lastPos:this.pos])
				this.lastPos = this.pos
			}
			for this.lastEmitPos < this.emitPos {
				fmt.Printf("%s, ", this.text[this.lastEmitPos])
				if this.text[this.lastEmitPos] < OP_ADJ {
					this.lastEmitPos++
					fmt.Printf(" %d\n", this.text[this.lastEmitPos])
				} else {
					fmt.Println()
				}
			}
			this.line++
		} else if ch == '#' {
			this.eatLine(config)
		} else if isIdentifierChar(ch) {
			this.parseIdentifier(config)
			return
		} else if isDigit(ch) {
			this.parseInt(config)
			return
		} else if ch == '/' {
			if this.src[this.pos] == '/' {
				this.eatLine(config)
			} else {
				this.token = DIV
			}
			return
		} else if (ch == '\'') || (ch == '"') {
			for ; this.pos < len(this.src); this.src[this.pos] != ch {
				//this.ival =
			}
		}
	}
}

func (this *Parser) eatLine(config *config.RunConfig) {
	for this.pos < len(this.src) && this.src[this.pos] != '\n' {
		this.pos++
	}
}

func (this *Parser) parseInt(config *config.RunConfig) {
	begin = this.pos - 1
	end = this.pos

	val := int(this.src[begin] - '0')

	if val != 0 {
		val := int(this.src[begin] - '0')
		for ; end < len(this.src); end++ {
			if !isDigit(src[end]) {
				break
			}
			val = 10*val + int(src[end]-'0')
		}
	} else if this.src[end] == 'x' || this.src[end] == 'x' {
		for ; end < len(this.src); end++ {
			if !isHex(src[end]) {
				break
			}
			val = 16*val + hexToInt(src[end])
		}
	} else {
		for ; end < len(this.src); end++ {
			if (src[end] < '0') || (src[end] > '7') {
				break
			}
			val = 8*val + int(src[end]-'0')
		}
	}

	this.pos = end
	this.ival = val
	this.token = NUM
}

func (this *Parser) parseIdentifier(config *config.RunConfig) {
	begin = this.pos - 1
	end = this.pos

	for end < len(this.src) {
		if !isIdentifierChar(src[end]) {
			break
		}
	}

	this.pos = end

	name := string(this.src[begin:end])

	symbol, ok := this.symbols[id]
	if ok {
		this.token = symbol.Token
		return
	}
	symbol = &Symbol{Token: ID, Name: name}
	this.symbols[id] = symbol

}

func hexToInt(ch byte) int {
	val := int(ch & 15)
	if ch >= 'A' {
		val += 9
	}
	return val
}

func isHex(ch byte) bool {
	return isDigit(ch) || isHexLower(ch) || isHexUpper(ch)
}

func isHexLower(ch byte) bool {
	return (ch >= 'a') && (ch <= 'f')
}

func isHexUpper(ch byte) bool {
	return (ch >= 'A') && (ch <= 'F')
}

func isAlphaNum(ch byte) bool {
	return isLower(ch) || isDigit(ch)
}

func isAlpha(ch byte) bool {
	return isLower(ch) || isUpper(ch)
}

func isLower(ch byte) {
	return (ch >= 'a') && (ch <= 'z')
}

func isUpper(ch byte) bool {
	return (ch >= 'A') && (ch <= 'Z')
}

func isDigit(ch byte) bool {
	return (ch >= '0') && (ch <= '9')
}

func isIdentifierFisrtChar(ch byte) bool {
	return isAlpha(ch) || (ch == '_')
}

func isIdentifierChar(ch byte) bool {
	return isAlphaNum(ch) || (ch == '_')
}
