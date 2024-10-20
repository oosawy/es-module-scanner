package scanner

import "unicode"

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
