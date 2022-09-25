package scanner

import (
	"devZ/lox/internal/tokens"
)

type Scanner struct {
	source []byte
	toks   []tokens.Token
}

func New(src []byte) (*Scanner, error) {
	return &Scanner{source: src}, nil
}

func (s *Scanner) ScanTokens() []tokens.Token {
	return s.toks
}
