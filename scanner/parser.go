package scanner

import "fmt"

func (s *scanner) tryParseImport() {
	start := s.position

	var tok Token
	if s.peekToken().Type != "IMPORT" {
		return
	} else {
		tok = s.readToken()
	}

	tok = s.readToken()
	switch tok.Type {
	case "STRING":
		goto speficer
	case "IDENT":
		tok = s.readToken()
		if tok.Type == "," {
			tok = s.readToken()
			goto namespaceOrNamed
		} else {
			goto fromClause
		}
	}

namespaceOrNamed:
	switch tok.Type {
	case "*":
		tok = s.readToken()
		if tok.Type != "AS" {
			s.syntaxError()
			return
		}

		tok = s.readToken()
		if tok.Type != "IDENT" {
			s.syntaxError()
			return
		}

	case "{":
		for {
			tok = s.readToken()
			println(tok.Type)
			if !tok.valid() {
				s.syntaxError()
				return
			}
			if tok.Type == "}" {
				break
			}
		}

	default:
		s.syntaxError()
		return
	}

	tok = s.readToken()
fromClause:
	if tok.Type != "FROM" {
		s.syntaxError()
		return
	}

	tok = s.readToken()
speficer:
	if tok.Type != "STRING" {
		s.syntaxError()
		return
	}

	if s.peekToken().Type == ";" {
		s.readToken()
	}
	lit := s.input[start:s.position]
	s.module.Imports = append(s.module.Imports, string(lit))

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
