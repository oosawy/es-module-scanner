package scanner

import "unicode"

type Span struct {
	Start int
	End   int
}

type Token struct {
	Type    string
	Literal string
	Span    Span
}

func (t Token) String() string {
	return t.Literal
}

func (t Token) valid() bool {
	return t.Type != "ILLIGAL" && t.Type != "EOF"
}

func (s *scanner) peekToken() Token {
	ch := s.ch
	position := s.position
	readPosition := s.readPosition
	defer func() {
		s.ch = ch
		s.position = position
		s.readPosition = readPosition
	}()

	return s.readToken()
}

func (s *scanner) readToken() Token {
	s.skipCommentWhiteSpace()

	start := s.position
	switch s.ch {
	case 0:
		return Token{Type: "EOF", Literal: "", Span: Span{start, s.position}}
	case '{':
		s.readChar()
		return Token{Type: "{", Literal: "{", Span: Span{start, s.position}}
	case '}':
		s.readChar()
		return Token{Type: "}", Literal: "}", Span: Span{start, s.position}}
	case ';':
		s.readChar()
		return Token{Type: ";", Literal: ";", Span: Span{start, s.position}}
	case ',':
		s.readChar()
		return Token{Type: ",", Literal: ",", Span: Span{start, s.position}}
	case '*':
		s.readChar()
		return Token{Type: "*", Literal: "*", Span: Span{start, s.position}}
	case '"', '\'':
		quote := s.ch

		for s.readChar() != 0 {
			if s.ch == quote {
				s.readChar()
				lit := string(s.input[start:s.position])
				return Token{Type: "STRING", Literal: lit, Span: Span{start, s.position}}
			} else if s.ch == '\\' {
				if s.readChar() == 0 {
					break
				}
			} else if isLineTerminator(s.ch) {
				break
			}
		}

		return Token{Type: "ILLIGAL", Literal: string(s.ch), Span: Span{start, s.position}}
	default:
		if l := findIdentifierStart(s.input[s.position:]); l > 0 {
			start := s.position

			for range l {
				s.readChar()
			}

			for {
				l := findIdentifierPart(s.input[s.position:])
				if l == 0 {
					break
				}
				for range l {
					s.readChar()
				}
			}

			lit := string(s.input[start:s.position])

			switch lit {
			case "import":
				return Token{Type: "IMPORT", Literal: lit, Span: Span{start, s.position}}
			case "export":
				return Token{Type: "EXPORT", Literal: lit, Span: Span{start, s.position}}
			case "from":
				return Token{Type: "FROM", Literal: lit, Span: Span{start, s.position}}
			case "as":
				return Token{Type: "AS", Literal: lit, Span: Span{start, s.position}}
			}

			return Token{Type: "IDENT", Literal: lit, Span: Span{start, s.position}}
		}

		return Token{Type: "ILLIGAL", Literal: string(s.ch), Span: Span{start, s.position}}
	}
}

func findIdentifierStart(chs []rune) int {
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

func findIdentifierPart(chs []rune) int {
	if l := findIdentifierStart(chs); l > 0 {
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
