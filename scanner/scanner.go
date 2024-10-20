package scanner

import (
	"fmt"
	"unicode"
)

type scanner struct {
	input        []rune
	position     int
	readPosition int
	ch           rune

	// contextDepth int

	module Module
}

type Module struct {
	Imports []string
	Exports []string
}

func Scan(input string) Module {
	s := &scanner{input: []rune(input)}

	for s.readChar() != 0 {
		s.skipCommentWhiteSpace()

		switch s.ch {
		case 'i':
			if s.isKeywordStart() {
				s.parseImport()
			}
		case 'e':
			if s.isKeywordStart() {
			}
		}
	}

	return s.module
}

func (s *scanner) syntaxError() {
	var line, col int
	for i, ch := range s.input {
		if i == s.position {
			break
		}
		if isLineTerminator(ch) {
			line++
			col = 0
		} else {
			col++
		}
	}

	var code, pointer string

	for i, ch := range s.input {
		if isLineTerminator(ch) {
			if i >= s.position {
				break
			}
			code = ""
			pointer = ""
		} else {
			code += string(ch)

			if i == s.position {
				pointer += "^"
			} else {
				pointer += " "
			}
		}
	}

	fmt.Printf("Syntax error at %d:%d (%d)\n\n%s\n%s\n", line, col, s.position, code, pointer)

	s.position = len(s.input)
}

func (s *scanner) parseImport() {
	position := s.position
	readPosition := s.readPosition

	tok := s.readToken()
	if tok.Type != "IMPORT" {
		goto bail
	}

	tok = s.readToken()
	switch tok.Type {
	case "STRING":
		s.skipCommentWhiteSpace()
		if s.ch == ';' {
			s.readChar()
		}
		lit := s.input[position:s.position]
		s.module.Imports = append(s.module.Imports, string(lit))
	case "{":
		for {
			tok = s.readToken()
			if tok.Type == "ILLIGAL" {
				s.syntaxError()
				return
			}
			if tok.Type == "}" {
				break
			}
		}

		tok = s.readToken()
		if tok.Type != "FROM" {
			goto bail
		}

		tok = s.readToken()
		if tok.Type != "STRING" {
			goto bail
		}

		s.skipCommentWhiteSpace()
		if s.ch == ';' {
			s.readChar()
		}

		lit := s.input[position:s.position]
		s.module.Imports = append(s.module.Imports, string(lit))
	}

	return

bail:
	s.position = position
	s.readPosition = readPosition
}

func (s *scanner) readChar() rune {
	if s.readPosition >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition += 1
	return s.ch
}

type Token struct {
	Type    string
	Literal string
}

func (s *scanner) readToken() Token {
	s.skipCommentWhiteSpace()

	switch s.ch {
	case 0:
		return Token{Type: "EOF", Literal: ""}
	case '{':
		s.readChar()
		return Token{Type: "{", Literal: "{"}
	case '}':
		s.readChar()
		return Token{Type: "}", Literal: "}"}
	case ';':
		s.readChar()
		return Token{Type: ";", Literal: ";"}
	case ',':
		s.readChar()
		return Token{Type: ",", Literal: ","}
	case '*':
		s.readChar()
		return Token{Type: "*", Literal: "*"}
	case '"', '\'':
		start := s.position
		quote := s.ch

		for s.readChar() != 0 {
			if s.ch == quote {
				s.readChar()
				return Token{Type: "STRING", Literal: string(s.input[start:s.position])}
			} else if s.ch == '\\' {
				if s.readChar() == 0 {
					break
				}
			} else if isLineTerminator(s.ch) {
				break
			}
		}

		return Token{Type: "ILLIGAL", Literal: string(s.ch)}
	default:
		if l := isIdentifierStart(s.input[s.position:]); l > 0 {
			start := s.position

			for range l {
				s.readChar()
			}

			for {
				l := isIdentifierPart(s.input[s.position:])
				if l == 0 {
					break
				}
				for range l {
					s.readChar()
				}
			}

			name := string(s.input[start:s.position])

			switch name {
			case "import":
				return Token{Type: "IMPORT", Literal: name}
			case "export":
				return Token{Type: "EXPORT", Literal: name}
			case "from":
				return Token{Type: "FROM", Literal: name}
			case "as":
				return Token{Type: "AS", Literal: name}
			}

			return Token{Type: "IDENT", Literal: name}
		}

		return Token{Type: "ILLIGAL", Literal: string(s.ch)}
	}
}

func (s *scanner) skipCommentWhiteSpace() {
	for {
		s.skipWhiteSpace()
		if s.ch == '/' {
			next := s.readChar()
			if next == '/' {
				for s.ch != 0 && !isLineTerminator(s.ch) {
					s.readChar()
				}
			} else if next == '*' {
				for {
					next = s.readChar()
					if next == 0 {
						break
					}
					if s.ch == '*' && next == '/' {
						s.readChar()
						break
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}
}

func (s *scanner) skipWhiteSpace() {
	for isWhiteSpace(s.ch) {
		if s.readChar() == 0 {
			break
		}
	}
}

func (s *scanner) isKeywordStart() bool {
	return s.position == 0 || isBoundary(s.input[s.position-1])
}

func isBoundary(ch rune) bool {
	return isWhiteSpace(ch) || isLineTerminator(ch) || (isPunctuator(ch) && ch != '.')
}

func isWhiteSpace(ch rune) bool {
	switch ch {
	case 0x09, 0x0B, 0x0C, 0xFEFF:
		return true
	}
	return unicode.Is(unicode.Zs, ch)
}

func isLineTerminator(ch rune) bool {
	switch ch {
	case 0x0A, 0x0D, 0x2028, 0x2029:
		return true
	}
	return false
}

func isPunctuator(ch rune) bool {
	switch ch {
	case '!', '%', '&', '(', ')', '*', '+', ',', '-', '.', '/', ':',
		';', '<', '=', '>', '?', '[', ']', '^', '{', '}', '|', '~':
		return true
	}
	return false
}

func isIdentifierStart(chs []rune) int {
	switch {
	case chs[0] == '$', chs[0] == '_',
		unicode.Is(unicode.L, chs[0]),
		unicode.Is(unicode.Nl, chs[0]),
		unicode.Is(unicode.Other_ID_Start, chs[0]):
		return 1
	case chs[0] == '\\':
		return readUnicodeEscapeSequence(chs[1:])
	}
	return 0
}

func isIdentifierPart(chs []rune) int {
	if l := isIdentifierStart(chs); l > 0 {
		return l
	}

	switch {
	case unicode.IsDigit(chs[0]),
		unicode.Is(unicode.Mn, chs[0]),
		unicode.Is(unicode.Mc, chs[0]),
		unicode.Is(unicode.Pc, chs[0]):
		return 1
	}
	return 0
}

func readUnicodeEscapeSequence(chs []rune) int {
	if chs[0] != 'u' {
		return 0
	}

	if chs[1] == '{' {
		for i := 2; i < len(chs); i++ {
			if chs[i] == '}' {
				return i + 1
			}
		}
		return 0
	} else {
		return 5
	}
}
