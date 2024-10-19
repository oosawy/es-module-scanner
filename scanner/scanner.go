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
		s.skipWhiteSpace()

		switch s.ch {
		case 'i':
			if s.isKeywordStart() && s.peekWord() == "import" {
				stmt := s.readLine()
				println(stmt)
				s.module.Imports = append(s.module.Imports, stmt)
			}
		case 'e':
			if s.isKeywordStart() && s.peekWord() == "export" {
				stmt := s.readLine()
				println(stmt)
				s.module.Exports = append(s.module.Exports, stmt)
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

func (s *scanner) peekWord() string {
	for p := s.readPosition; p < len(s.input); p++ {
		if isBoundary(s.input[p]) {
			return string(s.input[s.position:p])
		}
	}
	return string(s.input[s.position:])
}

func (s *scanner) readLine() string {
	start := s.position
	for !isLineTerminator(s.ch) {
		if s.readChar() == 0 {
			break
		}
	}
	return string(s.input[start:s.position])
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
