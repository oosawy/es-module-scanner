package scanner

import (
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
				s.tryParseImport()
			}
		case 'e':
			if s.isKeywordStart() {
			}
		}
	}

	return s.module
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

func (s *scanner) peekChar() rune {
	if s.readPosition >= len(s.input) {
		return 0
	}
	return s.input[s.readPosition]
}

func (s *scanner) skipCommentWhiteSpace() {
	for {
		s.skipWhiteSpace()
		if s.ch == '/' {
			switch s.peekChar() {
			case '/':
				s.readChar()
				for s.readChar() != 0 {
					if isLineTerminator(s.ch) {
						break
					}
				}

			case '*':
				s.readChar()
				for s.readChar() != 0 {
					if s.ch == '*' && s.peekChar() == '/' {
						s.readChar()
						s.readChar()
						break
					}
				}

			default:
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
