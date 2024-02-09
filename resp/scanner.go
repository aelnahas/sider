package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

var eof = rune(0)

const (
	invalidSize = -1
	unsetSize   = -2
)

type Scanner struct {
	r     *bufio.Reader
	count int
	size  int
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r), count: unsetSize}
}

func (s *Scanner) HasNext() bool {
	return s.count == unsetSize || s.count > 0
}

func (s *Scanner) Next() (Token, string, error) {
	if s.count == unsetSize {
		count, err := s.readArray()
		if err != nil {
			return TokenEOF, string(eof), err
		}
		s.count = count
		s.size = count
	}

	if s.count <= 0 {
		return TokenEOF, string(eof), ErrOutOfBound
	}

	token, lit, err := s.Scan()
	if err == nil {
		s.count--
	}
	return token, lit, err
}

func (s *Scanner) readArray() (int, error) {
	ch, err := s.read()
	if err != nil {
		return -1, err
	}

	if ch != rune(SymbolArray) {
		return -1, ErrNotABulkString
	}

	size, err := s.readSize()
	if err != nil {
		return -1, err
	}

	return size, err
}

func (s *Scanner) Scan() (Token, string, error) {
	ch, err := s.read()
	if err != nil {
		return TokenEOF, string(eof), err
	}

	if s.isCRLF(ch) {
		if err := s.unread(); err != nil {
			return TokenEOF, string(eof), err
		}

		if err := s.readCRLF(); err != nil {
			return TokenEOF, string(eof), err
		}

		return s.Scan()
	}

	if ch != rune(SymbolBulkString) {
		return TokenEOF, string(eof), ErrNotABulkString
	}

	word, err := s.readWord()
	if err != nil {
		return TokenEOF, string(eof), err
	}

	err = s.readCRLF()
	if err != nil {
		return TokenEOF, string(eof), err
	}

	// we are scanning the first item which should be the the command
	if s.size == s.count {
		return s.getCommandToken(word)
	}

	return TokenArg, word, nil
}

func (s *Scanner) getCommandToken(word string) (Token, string, error) {

	switch strings.ToUpper(word) {
	case CmdGet:
		return TokenGet, word, nil
	case CmdPing:
		return TokenPing, word, nil
	case CmdSet:
		return TokenSet, word, nil
	case CmdEcho:
		return TokenEcho, word, nil
	case CmdExists:
		return TokenExists, word, nil
	case CmdDel:
		return TokenDel, word, nil
	case CmdPub:
		return TokenPub, word, nil
	case CmdSub:
		return TokenSub, word, nil
	case CmdUnSub:
		return TokenUnSub, word, nil
	default:
		return TokenEOF, word, ErrUnknownCommand{Name: strings.ToUpper(word)}
	}
}

func (s *Scanner) readWord() (string, error) {
	var buf bytes.Buffer

	size, err := s.readSize()
	if err != nil {
		return "", err
	}

	for i := 0; i < size; i++ {
		ch, err := s.read()
		if err != nil {
			return "", err
		}

		_, err = buf.WriteRune(ch)
		if err != nil {
			return "", err
		}
	}

	return buf.String(), nil
}

func (s *Scanner) readSize() (int, error) {
	size := 0

	for {
		ch, err := s.read()
		if err != nil {
			return -1, err
		}

		if unicode.IsDigit(ch) {
			size = size*10 + int(ch-'0')
		} else if ch == rune(SymbolCR) {
			if err := s.unread(); err != nil {
				return -1, err
			}
			break
		} else {
			return -1, fmt.Errorf("expected digits only when reading array or bulk string size, current char (%c)", ch)
		}

	}

	if err := s.readCRLF(); err != nil {
		return -1, err
	}

	return size, nil
}

func (s *Scanner) readCRLF() error {
	ch, err := s.read()
	if err != nil {
		return err
	}
	if ch != rune(SymbolCR) {
		return ErrUnexpectedSymbol{Wanted: rune(SymbolCR), Got: ch}
	}

	ch, err = s.read()
	if err != nil {
		return err
	}
	if ch != rune(SymbolLF) {
		return ErrUnexpectedSymbol{Wanted: rune(SymbolLF), Got: ch}
	}

	return nil
}

func (s *Scanner) isCRLF(ch rune) bool {
	return ch == rune(SymbolCR) || ch == rune(SymbolLF)
}

func (s *Scanner) read() (rune, error) {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof, err
	}

	return ch, nil
}

func (s *Scanner) unread() error {
	return s.r.UnreadRune()
}
