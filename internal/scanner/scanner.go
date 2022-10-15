package scanner

import (
	"bytes"
	"strconv"
	"unicode/utf8"

	"devZ/lox/internal/reporters"
	"devZ/lox/internal/tokens"
)

var newLines = []byte{10, 11, 12, 13, 133}

func IsWesternDigit(r rune) bool {
	return r >= 48 && r <= 57
}

func IsEnglishAlpha(r rune) bool {
	lower := r >= 97 && r <= 122
	upper := r >= 65 && r <= 90
	under := r == 95
	return lower || upper || under
}

func IsAlphaNumeric(r rune) bool {
	return IsEnglishAlpha(r) || IsWesternDigit(r)
}

type Scanner struct {
	source  []byte
	toks    []*tokens.Token
	srcLn   int
	curIdx  int
	tokSt   int
	accum   *reporters.Accumulator
	printer *reporters.PrettyPrinter
}

func New(src []byte, a *reporters.Accumulator) *Scanner {
	return &Scanner{source: src, accum: a}
}

func (s *Scanner) AddPrinter(pp *reporters.PrettyPrinter) {
	s.printer = pp
}

func (s *Scanner) lexeme() string {
	return string(s.source[s.tokSt:s.curIdx])
}

func (s *Scanner) advance() rune {
	r, sz := utf8.DecodeRune(s.source[s.curIdx:])
	s.curIdx += sz
	return r
}

func (s *Scanner) match(check rune) bool {
	if s.isAtEnd() {
		return false
	}
	r, sz := utf8.DecodeRune(s.source[s.curIdx:])
	if r != check {
		return false
	}
	// advance only if next rune matches check rune
	s.curIdx += sz
	return true
}

func (s *Scanner) peek() rune {
	r, _ := utf8.DecodeRune(s.source[s.curIdx:])
	return r
}

func (s *Scanner) peekNext() rune {
	_, sz := utf8.DecodeRune(s.source[s.curIdx:])
	r, _ := utf8.DecodeRune(s.source[s.curIdx+sz:])
	return r
}

func (s *Scanner) string() {
	// look for the terminating quote
	for s.peek() != '"' && !s.isAtEnd() {
		if bytes.ContainsRune(newLines, s.peek()) {
			s.srcLn++
			s.advance()
		}
	}
	// if you hit the end of the source, something's missing
	if s.isAtEnd() {
		s.accum.AddError(s.srcLn, "Unterminated string.")
	}
	// consume the terminating quote
	s.advance()
	// get the string literal value
	value := string(s.source[s.tokSt+1 : s.curIdx-1])
	s.addToken(tokens.STRING, value)
}

func (s *Scanner) number() {
	for IsWesternDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && IsWesternDigit(s.peekNext()) {
		s.advance()
		for IsWesternDigit(s.peek()) {
			s.advance()
		}
	}
	val, err := strconv.ParseFloat(s.lexeme(), 64)
	if err != nil {
		panic(err)
	}
	s.addToken(tokens.NUMBER, val)
}

func (s *Scanner) identifier() {
	for IsAlphaNumeric(s.peek()) {
		s.advance()
	}
	val := s.lexeme()
	tid, prs := tokens.Keywords[val]
	if !prs {
		tid = tokens.IDENTIFIER
	}
	s.addToken(tid, nil)
}

func (s *Scanner) checkEquals(short, long tokens.TokID) {
	if s.match('=') {
		s.addToken(long, nil)
	} else {
		s.addToken(short, nil)
	}
}

func (s *Scanner) checkComment() {
	if s.match('/') {
		var chars []rune
		for s.peek() != 10 && !s.isAtEnd() {
			chars = append(chars, s.advance())
		}
		s.accum.AddComment(s.srcLn, string(chars))
	} else {
		s.addToken(tokens.SLASH, nil)
	}
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
	// consume and store the next rune of source
	c := s.advance()
	// tokens are detected by their first rune
	switch c {
	// grouping & closures
	case '(':
		s.addToken(tokens.LEFT_PAREN, nil)
	case ')':
		s.addToken(tokens.RIGHT_PAREN, nil)
	case '{':
		s.addToken(tokens.RIGHT_BRACE, nil)
	case '}':
		s.addToken(tokens.LEFT_BRACE, nil)
	// separators & line ends
	case ',':
		s.addToken(tokens.COMMA, nil)
	case '.':
		s.addToken(tokens.DOT, nil)
	case ';':
		s.addToken(tokens.SEMICOLON, nil)
	// arithmetic
	case '-':
		s.addToken(tokens.MINUS, nil)
	case '+':
		s.addToken(tokens.PLUS, nil)
	case '*':
		s.addToken(tokens.STAR, nil)
	// slash is for division, but comments start with slashes (leo_pointing.jpg)
	case '/':
		s.checkComment()
	// negation, assignment, equality & comparison operators
	case '!':
		s.checkEquals(tokens.BANG, tokens.BANG_EQUAL)
	case '=':
		s.checkEquals(tokens.EQUAL, tokens.EQUAL_EQUAL)
	case '>':
		s.checkEquals(tokens.GREATER, tokens.GREATER_EQUAL)
	case '<':
		s.checkEquals(tokens.LESS, tokens.LESS_EQUAL)
	// strings start with quotes, as is their wont
	case '"':
		s.string()
	// ignore non linebreak whitespace
	case 9, 32, 160:
	// linebreaks
	case 10, 11, 12, 13, 133:
		s.srcLn++
	// after everything else, check for numbers or identifiers; if not a number
	// or identifier, add an error and continue scanning
	default:
		if IsWesternDigit(c) {
			s.number()
		} else if IsEnglishAlpha(c) {
			s.identifier()
		} else {
			s.accum.AddError(s.srcLn, "Unexpected character.")
		}
	}
}

func (s *Scanner) ScanTokens() []*tokens.Token {
	for !s.isAtEnd() {
		s.tokSt = s.curIdx
		s.scanToken()
	}
	s.addToken(tokens.EOF, nil)
	return s.toks
}

func (s *Scanner) CurLine() int {
	return s.srcLn
}
