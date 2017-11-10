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
	Class  Token
	Type   Type
	Value  interface{}
	BClass Token
	BType  Type
	BValue interface{}
}

type Symbols map[string]*Symbol

type Parser struct {
	src             []byte
	pos             uint64
	lastPos         uint64
	emitPos         uint64
	lastEmitPos     uint64
	textPos         uint64
	dataPos         uint64
	vim             *vm.Vim
	line            uint64
	currentExprType uint64
	token           Token
	tokenValue      uint64
	symbols         Symbols
	ival            uint64
	sval            []byte
	currentId       *Symbol
	main            *Symbol
	indexOfBP       uint64
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

func (this *Parser) CurrentIdName() string {
	return this.currentId.Name
}

func (this *Parser) InitVim(config *config.RunConfig, textSize, dataSize, stackSize uint64) {
	this.vim = vm.NewVM(textSize, dataSize, stackSize)
}

func (this *Parser) InitKeywords(config *config.RunConfig) bool {
	keywords := []string{"char", "int", "void", "if", "else", "while", "return", "enum", "sizeof", "printf", "main"}
	this.dataPos = this.vim.DataBegin()
	for _, v := range keywords {
		symbol := &Symbol{Token: TOKEN_ID, Name: v, Class: TOKEN_SYS, Type: TYPE_INT}
		this.symbols[v] = symbol
		ok, pos := this.vim.AddDataString(this.dataPos, v)
		if !ok {
			fmt.Printf("ERROR: init keyword \"%s\" failed\n", v)
			return false
		}
		this.dataPos = pos
	}

	this.symbols["void"].Token = TOKEN_CHAR
	this.main, _ = this.symbols["main"]
	return true
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

func (this *Parser) global_declaration(config *config.RunConfig) bool {
	// global_declaration ::= enum_decl | variable_decl | function_decl
	//
	// enum_decl ::= 'enum' [id] '{' id ['=' 'num'] {',' id ['=' 'num'} '}'
	//
	// variable_decl ::= type {'*'} id { ',' {'*'} id } ';'
	//
	// function_decl ::= type {'*'} id '(' parameter_decl ')' '{' body_decl '}'
	baseType := TYPE_INT
	currentType := TYPE_INT

	// parse enum, this should be treated alone.
	if this.token == TOKEN_ENUM {
		return this.enum_declaration(config)
	}

	// parse type information
	if this.token == TOKEN_INT {
		this.Next(config)
	} else if this.token == TOKEN_CHAR {
		this.Next(config)
		baseType = TYPE_CHAR
	}

	// parse the comma seperated variable declaration.
	for this.token != ';' && this.token != '}' {
		currentType = baseType

		// parse pointer type, note that there may exist `int ****x;`
		for this.token == TOKEN_MUL {
			this.Next(config)
			currentType += TYPE_PTR
		}

		if this.token != TOKEN_ID {
			// invalid declaration
			fmt.Printf("ERROR: line %d, bad global declaration\n", this.line)
			return false
		}

		if this.currentId.Class > 0 {
			// identifier exists
			fmt.Printf("ERROR: line %d, duplicate global declaration\n")
			return false
		}

		this.Next(config)
		this.currentId.Type = currentType

		if this.token == '(' {
			this.currentId.Class = TOKEN_FUN
			this.currentId.Value = this.textPos
			if !this.function_declaration(config) {
				return false
			}
		} else {
			// variable declaration
			this.currentId.Class = TOKEN_GLO
			this.currentId.Value = this.dataPos
			this.dataPos += vm.SizeofUint64()
		}

		if this.token == ',' {
			this.Next(config)
		}
	}

	return true
}

func (this *Parser) function_declaration(config *config.RunConfig) bool {
	// function_decl ::= type {'*'} id '(' parameter_decl ')' '{' body_decl '}'
	this.Next(config)

	if !this.function_parameter(config) {
		return false
	}

	if this.token != ')' {
		fmt.Printf("ERROR: line %d, need ')'\n")
		return false
	}
	this.Next(config)

	if this.token != '{' {
		fmt.Printf("ERROR: line %d, need '{'\n")
		return false
	}
	this.Next(config)

	if !this.function_body(config) {
		return false
	}

	// unwind local variable declarations for all local variables.
	this.unwindLocal(config)

	return true
}

func (this *Parser) unwindLocal(config *config.RunConfig) {
	for _, v := range this.symbols {

		if v.Class == TOKEN_LOC {
			v.Class = v.BClass
			v.Type = v.BType
			v.Value = v.BValue
		}
	}
}

func (this *Parser) function_parameter(config *config.RunConfig) bool {
	// parameter_decl ::= type {'*'} id {',' type {'*'} id}
	paramNum := 0
	for this.token != ')' {
		paramType := TYPE_INT

		// int name, ...
		if this.token == TOKEN_INT {
			this.Next(config)
		} else if this.token == TOKEN_CHAR {
			paramType = TYPE_CHAR
			this.Next(config)
		} else {
			fmt.Printf("ERROR: line %d, no parameter type\n")
			return false
		}

		// pointer type
		for this.token == TOKEN_MUL {
			this.Next(config)
			paramType += TYPE_PTR
		}

		// parameter name
		if this.token != TOKEN_ID {
			fmt.Printf("ERROR: line %d, no parameter name\n")
			return false
		}

		if this.currentId.Class == TOKEN_LOC {
			fmt.Printf("ERROR: line %d, duplicate parameter declaration\n")
			return false
		}

		this.Next(config)

		// store the local variable
		this.storeLocal(config)
		this.currentId.Type = paramType
		this.currentId.Value = paramNum // index of current parameter
		paramNum++

		if this.token == ',' {
			this.Next(config)
		}
	}

	this.indexOfBP = uint64(paramNum) + 1
	return true
}

func (this *Parser) function_body(config *config.RunConfig) bool {
	// type func_name (...) {...}
	//                   -->|   |<--
	// ... {
	// 1. local declarations
	// 2. statements
	// }

	posLocal := this.indexOfBP

	for this.token == TOKEN_INT || this.token == TOKEN_CHAR {
		baseType := TYPE_INT
		if this.token == TOKEN_CHAR {
			baseType = TYPE_CHAR
		}

		for this.token != ';' {
			varType := baseType
			for this.token == TOKEN_MUL {
				this.Next(config)
				varType += TYPE_PTR
			}

			if this.token != TOKEN_ID {
				fmt.Printf("ERROR: line %d, no local var name\n")
				return false
			}

			if this.currentId.Class == TOKEN_LOC {
				fmt.Printf("ERROR: line %d, duplicate parameter declaration\n")
				return false
			}

			this.Next(config)

			// store the local variable
			this.storeLocal(config)
			this.currentId.Type = varType
			this.currentId.Value = posLocal
			posLocal++

			if this.token == ',' {
				this.Next(config)
			}
		}

		this.Next(config)
	}

	// save the stack size for local variables
	var ok bool
	ok, this.textPos = this.vim.AddOpCode(this.textPos, vm.OP_ENT)
	if !ok {
		fmt.Printf("ERROR: line %d, generate OP_ENT failed\n")
		return false
	}

	ok, this.textPos = this.vim.AddUint64(this.textPos, posLocal-this.indexOfBP)
	if !ok {
		fmt.Printf("ERROR: line %d, generate OP_ENT param failed\n")
		return false
	}

	// statements
	if !this.statement(config) {
		return false
	}

	ok, this.textPos = this.vim.AddOpCode(this.textPos, vm.OP_LEV)
	if !ok {
		fmt.Printf("ERROR: line %d, generate OP_LEV failed\n")
		return false
	}

	return true
}

func (this *Parser) statement(config *config.RunConfig) bool {
	return false
}

func (this *Parser) storeLocal(config *config.RunConfig) {
	this.currentId.BClass = this.currentId.Class
	this.currentId.Class = TOKEN_LOC
	this.currentId.BType = this.currentId.Type
	this.currentId.BValue = this.currentId.Value
}

func (this *Parser) enum_declaration(config *config.RunConfig) bool {
	// parse enum [id] { a = 1, b = 3, ...}
	if this.token != TOKEN_ENUM {
		return false
	}

	this.Next(config)
	if this.token != '{' {
		if this.token != TOKEN_ID {
			fmt.Printf("ERROR: line %d, unexpected token %s\n", this.line, this.token)
			return false
		}
		// skip enum name
		this.Next(config)
	}
	if this.token != '{' {
		fmt.Printf("ERROR: line %d, no '{'\n", this.line)
		return false
	}
	this.Next(config)

	val := uint64(0)
	for this.token != '}' {
		if this.token != TOKEN_ID {
			fmt.Printf("ERROR: line %d, bad enum identifier %s\n", this.line, this.token)
			return false
		}
		this.Next(config)
		if this.token == TOKEN_ASSIGN {
			this.Next(config)
			if this.token != TOKEN_NUM {
				fmt.Printf("ERROR: line %d, bad enum initializer\n", this.line)
				return false
			}
			val = this.ival
			this.Next(config)
		}

		this.currentId.Class = TOKEN_NUM
		this.currentId.Type = TYPE_INT
		this.currentId.Value = val
		val++

		if this.token == ',' {
			this.Next(config)
		}
	}

	if this.token != '}' {
		fmt.Printf("ERROR: line %d, no '}'\n", this.line, this.token)
		return false
	}
	this.Next(config)

	if this.token != ';' {
		fmt.Printf("ERROR: line %d, no ';'\n", this.line, this.token)
		return false
	}
	this.Next(config)

	return true
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
		this.currentId = symbol
		return
	}
	this.token = TOKEN_ID
	symbol = &Symbol{Token: TOKEN_ID, Name: name}
	this.symbols[name] = symbol
	this.currentId = symbol

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
