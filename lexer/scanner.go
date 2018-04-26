package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

const bufferSize = 32

type Scanner struct {
	rd    *bufio.Reader
	Token *Token
}

func NewScanner(rd io.Reader) *Scanner {
	return &Scanner{rd: bufio.NewReader(rd)}
}

func (s *Scanner) ReadNext() error {
	err := s.readNext()
	if err == io.EOF {
		s.Token = &Token{Type: EOF, Value: "EOF"}
		return nil
	}
	return err
}

func (s *Scanner) readNext() error {
	var r rune
	var err error

	for {
		// Skip whitespace

		r, _, err = s.rd.ReadRune()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(r) {
			break
		}
	}

	if r == '"' {
		// String
		v, err := scanString(s.rd)
		if err != nil {
			return err
		}
		s.Token = &Token{Type: String, Value: v}
		return nil
	}

	if unicode.IsLetter(r) || r == '_' {
		// Identifier
		v, err := scanIdentifier(s.rd)
		if err != nil {
			return err
		}
		id := string(r) + v
		s.Token = &Token{Type: Identifier, Value: id}
		return nil
	}

	if unicode.IsDigit(r) {
		// Number
		v, err := scanIdentifier(s.rd)
		if err != nil {
			return err
		}
		num := string(r) + v
		s.Token = &Token{Type: Integer, Value: num}
		return nil
	}

	if r == ':' {
		r, _, err = s.rd.ReadRune()
		if err != nil {
			return err
		}
		if r != '=' {
			return fmt.Errorf("Unexpected symbol %c", r)
		}
		s.Token = &Token{Type: DeclarationOperator, Value: ":="}
		return nil
	}

	if r == '=' {
		s.Token = &Token{Type: AssignmentOperator, Value: "="}
		return nil
	}

	if r == ';' {
		s.Token = &Token{Type: Semicolon, Value: ";"}
		return nil
	}

	if r == '(' {
		s.Token = &Token{Type: LeftParen, Value: "("}
		return nil
	}

	if r == ')' {
		s.Token = &Token{Type: RightParen, Value: ")"}
		return nil
	}

	if r == '*' {
		s.Token = &Token{Type: MultiplicationOperator, Value: "*"}
		return nil
	}

	if r == '/' {
		s.Token = &Token{Type: DivisionOperator, Value: "/"}
		return nil
	}

	if r == '+' {
		s.Token = &Token{Type: AdditionOperator, Value: "+"}
		return nil
	}

	if r == '-' {
		s.Token = &Token{Type: SubstractionOperator, Value: "-"}
		return nil
	}

	return fmt.Errorf("Unexpected symbol %c", r)
}

func scanString(rd *bufio.Reader) (string, error) {
	s := ""
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			return "", err
		}

		if r == '"' {
			return s, nil
		}

		s += string(r)
	}
}

func scanIdentifier(rd *bufio.Reader) (string, error) {
	s := ""
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			return "", err
		}

		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			if err := rd.UnreadRune(); err != nil {
				return "", err
			}
			return s, nil
		}

		s += string(r)
	}
}

func scanInteger(rd *bufio.Reader) (string, error) {
	s := ""
	for {
		r, _, err := rd.ReadRune()
		if err != nil {
			return "", err
		}

		if !unicode.IsDigit(r) {
			if err := rd.UnreadRune(); err != nil {
				return "", err
			}
			return s, nil
		}

		s += string(r)
	}
}