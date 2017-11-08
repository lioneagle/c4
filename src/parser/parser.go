package parser

import (
	"config"
	"fmt"
	"io/ioutil"
	"strconv"
	//"strings"
	"vm"
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

type Symbols map[string]*Symbol

type Parser struct {
	src             []byte
	pos             uint64
	lastPos         uint64
	emitPos         uint64
	lastEmitPos     uint64
	vim             vm.Vim
	line            uint64
	currentExprType uint64
	token           Token
	tokenValue      uint64
	symbols         Symbols
	ival            uint64
	sval            []byte
	idName          []byte
}

func NewParser() *Parser {
	return &Parser{symbols: make(map[string]*Symbol)}
}

func (this *Parser) ReadFile(filename string) (err error) {
	this.src, err = ioutil.ReadFile(filename)
	if err == nil && len(this.src) > 0 {
		this.line = 1
	}
	return err
}

func (this *Parser) Line() uint64 {
	return this.line
}

func (this *Parser) Eof() bool {
	return this.token == TOKEN_EOF
}

func (this *Parser) Token() Token {
	return this.token
}

func (this *Parser) IdName() string {
	return string(this.idName)
}

func (this *Parser) Next(config *config.RunConfig) {
	for this.pos < uint64(len(this.src)) {
		ch := this.src[this.pos]
		//fmt.Printf("Next: ch = %c, ch = %d, line = %d\n", ch, ch, this.line)
		//fmt.Printf("Next: ch = %c, ch = %d\n", ch, int(ch))
		//fmt.Printf("Next: ch = %s\n", PrintChar(ch))
		//fmt.Println("here1")
		this.pos++

		if ch == '\n' {
			if config.PrintSource {
				fmt.Printf("%d: %s", this.line, this.src[this.lastPos:this.pos])
				this.lastPos = this.pos
			}
			for this.lastEmitPos < this.emitPos {
				op := this.vim.OpCode(this.lastEmitPos)
				fmt.Printf("%s, ", op)
				if op < vm.OP_ADJ {
					this.lastEmitPos++
					fmt.Printf(" %d\n", this.vim.Uint64(this.lastEmitPos))
				} else {
					fmt.Println()
				}
			}
			this.line++
			//fmt.Printf("line = %d\n", this.line)
		} else if ch == '#' {
			this.eatLine(config)
		} else if isIdentifierFisrtChar(ch) {
			this.parseIdentifier(config)
			return
		} else if isDigit(ch) {
			this.parseInt(config)
			return
		} else if ch == '/' {
			if this.src[this.pos] == '/' {
				this.eatLine(config)
			} else {
				this.token = TOKEN_DIV
			}
			return
		} else if ch == '\'' {
			this.parseChar(config)
		} else if ch == '"' {
			this.parseString(config)
		} else if ch == '=' {
			if this.src[this.pos] == '=' {
				this.token = TOKEN_EQ
			} else {
				this.token = TOKEN_ASSIGN
			}
			return
		} else if ch == '+' {
			if this.src[this.pos] == '+' {
				this.token = TOKEN_INC
			} else {
				this.token = TOKEN_ADD
			}
			return
		} else if ch == '-' {
			if this.src[this.pos] == '-' {
				this.token = TOKEN_DEC
			} else {
				this.token = TOKEN_SUB
			}
			return
		} else if ch == '!' {
			if this.src[this.pos] == '=' {
				this.token = TOKEN_NE
			}
			return
		} else if ch == '<' {
			if this.src[this.pos] == '=' {
				this.token = TOKEN_LE
			} else if this.src[this.pos] == '<' {
				this.token = TOKEN_SHL
			} else {
				this.token = TOKEN_LT
			}
			return
		} else if ch == '>' {
			if this.src[this.pos] == '=' {
				this.token = TOKEN_GE
			} else if this.src[this.pos] == '>' {
				this.token = TOKEN_SHR
			} else {
				this.token = TOKEN_GT
			}
			return
		} else if ch == '|' {
			if this.src[this.pos] == '|' {
				this.token = TOKEN_LOR
			} else {
				this.token = TOKEN_OR
			}
			return
		} else if ch == '&' {
			if this.src[this.pos] == '&' {
				this.token = TOKEN_LAN
			} else {
				this.token = TOKEN_AND
			}
			return
		} else if ch == '^' {
			this.token = TOKEN_XOR
			return
		} else if ch == '%' {
			this.token = TOKEN_MOD
			return
		} else if ch == '*' {
			this.token = TOKEN_MUL
			return
		} else if ch == '[' {
			this.token = TOKEN_BRAK
			return
		} else if ch == '?' {
			this.token = TOKEN_COND
			return
		} else if ch == '^' {
			this.token = TOKEN_XOR
			return
		} else if ch == '~' || ch == ';' || ch == '{' || ch == '}' || ch == '(' || ch == ')' || ch == ']' || ch == ',' || ch == ':' {
			this.token = Token(ch)
			return
		}
	}
	this.token = TOKEN_EOF
}

func (this *Parser) Expr(config *config.RunConfig) bool {
	if this.token == TOKEN_EOF {
		this.Error("unexpected eof in expression")
		return false
	}

	if this.token == TOKEN_NUM {

	}

	return true

}

func (this *Parser) Error(format string, args ...interface{}) {
	fmt.Printf("%d: %s\n", this.line, fmt.Sprintf(format, args...))
}

func (this *Parser) parseString(config *config.RunConfig) {
	begin := this.pos
	for this.pos < uint64(len(this.src)) && this.src[this.pos] != '"' {
		this.pos++
	}
	if this.pos >= uint64(len(this.src)) {
		fmt.Printf("ERROR: invalid string at pos %d\n", this.pos)
		return
	}

	this.token = TOKEN_STRING
	this.sval = this.src[begin:this.pos]
	this.pos++
}

func (this *Parser) parseChar(config *config.RunConfig) {
	end := this.pos
	if end >= uint64(len(this.src)) {
		return
	}

	if this.src[end] != '\\' {
		this.ival = uint64(this.src[end])
	} else {
		end++
		if end > uint64(len(this.src)) {
			return
		}
		switch this.src[end] {
		case 'n':
			this.ival = '\n'
		case 'r':
			this.ival = '\r'
		case 't':
			this.ival = '\t'
		default:
			fmt.Printf("ERROR: unknown char at pos %d\n", end)
			return
		}
		end++

		if this.src[end] != '\'' {
			fmt.Printf("ERROR: invalid char at pos %d\n", end)
			return
		}
	}
	this.token = TOKEN_CHAR
}

func (this *Parser) eatLine(config *config.RunConfig) {
	for this.pos < uint64(len(this.src)) && this.src[this.pos] != '\n' {
		//fmt.Printf("Next: ch = %s\n", PrintChar(this.src[this.pos]))
		this.pos++
	}
	//fmt.Printf("eatLine: end, line = %d\n", this.line)
}

func (this *Parser) parseInt(config *config.RunConfig) {
	begin := this.pos - 1
	end := this.pos

	val := uint64(this.src[begin] - '0')

	if val != 0 {
		val := int(this.src[begin] - '0')
		for ; end < uint64(len(this.src)); end++ {
			if !isDigit(this.src[end]) {
				break
			}
			val = 10*val + int(this.src[end]-'0')
		}
	} else if this.src[end] == 'x' || this.src[end] == 'x' {
		for ; end < uint64(len(this.src)); end++ {
			if !isHex(this.src[end]) {
				break
			}
			val = 16*val + uint64(hexToInt(this.src[end]))
		}
	} else {
		for ; end < uint64(len(this.src)); end++ {
			if (this.src[end] < '0') || (this.src[end] > '7') {
				break
			}
			val = 8*val + uint64(this.src[end]-'0')
		}
	}

	this.pos = end
	this.ival = val
	this.token = TOKEN_NUM
}

func (this *Parser) parseIdentifier(config *config.RunConfig) {
	begin := this.pos - 1
	end := this.pos

	for ; end < uint64(len(this.src)); end++ {
		if !isIdentifierChar(this.src[end]) {
			break
		}
	}

	this.pos = end

	name := string(this.src[begin:end])

	//fmt.Println("parseIdentifier: name =", name)

	symbol, ok := this.symbols[name]
	if ok {
		//fmt.Println("parseIdentifier: token =", symbol.Token)
		this.token = symbol.Token
		this.idName = this.src[begin:end]
		return
	}
	this.token = TOKEN_ID
	symbol = &Symbol{Token: TOKEN_ID, Name: name}
	this.symbols[name] = symbol
	this.idName = this.src[begin:end]

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

func isLower(ch byte) bool {
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

func PrintChar(ch byte) string {
	if strconv.IsPrint(rune(ch)) {
		return fmt.Sprintf("%c", ch)
	}
	return fmt.Sprintf("%d", ch)

}
