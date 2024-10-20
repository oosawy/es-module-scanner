package scanner

import "fmt"

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
