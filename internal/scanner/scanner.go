package scanner

import (
	"devZ/lox/internal/tokens"
	"unicode/utf8"
)

type Scanner struct {
	source []byte
	toks   []*tokens.Token
	srcLn  int
	curIdx int
	tokSt  int
}

func New(src []byte) (*Scanner, error) {
	return &Scanner{source: src}, nil
}

func (s *Scanner) lexeme() string {
	return string(s.source[s.tokSt:s.curIdx])
}

func (s *Scanner) advance() rune {
	r, sz := utf8.DecodeRune(s.source[s.curIdx:])
	s.curIdx += sz
	return r
}

func (s *Scanner) isAtEnd() bool {
	return s.curIdx >= len(s.source)
}

func (s *Scanner) addToken(tid tokens.TokID, obj interface{}) {
	tok, err := tokens.New(tid, s.lexeme(), s.srcLn, obj)
	if err != nil {
		panic(err)
	}
	s.toks = append(s.toks, tok)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case utf8.RuneError:
		s.addToken(tokens.EOF, nil)
	case '(':
		s.addToken(tokens.RIGHT_PAREN, nil)
	case ')':
		s.addToken(tokens.LEFT_PAREN, nil)
	case '{':
		s.addToken(tokens.RIGHT_BRACE, nil)
	case '}':
		s.addToken(tokens.LEFT_BRACE, nil)
	case ',':
		s.addToken(tokens.COMMA, nil)
	case '.':
		s.addToken(tokens.DOT, nil)
	case '-':
		s.addToken(tokens.MINUS, nil)
	case '+':
		s.addToken(tokens.PLUS, nil)
	case ';':
		s.addToken(tokens.SEMICOLON, nil)
	case '*':
		s.addToken(tokens.STAR, nil)
	}
}

func (s *Scanner) ScanTokens() []*tokens.Token {
	for !s.isAtEnd() {
		s.tokSt = s.curIdx
		s.scanToken()
	}
	end, err := tokens.New(tokens.EOF, "", s.srcLn, nil)
	if err != nil {
		panic(err)
	}
	s.addToken(end)
	return s.toks
}

func (s *Scanner) CurLine() int {
	return s.srcLn
}
