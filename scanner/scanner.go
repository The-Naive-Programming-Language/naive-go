package scanner

import (
	"fmt"
	"unicode/utf8"

	"naive/token"
)

const (
	eof = -1
	bom = 0xFEFF
)

type Scanner struct {
	// immutable
	file *token.File
	src  []byte

	ch         rune // current character
	offset     int  // character offset
	rdOffset   int  // reading offset (start position of next character)
	lineOffset int

	NumErrors int
}

func New(file *token.File, src []byte) *Scanner {
	s := &Scanner{
		file: file,
		src:  src,

		// whitespace does not matter
		ch:         ' ',
		offset:     0,
		rdOffset:   0,
		lineOffset: 0,
	}

	s.next()
	if s.ch == bom {
		s.next()
	}

	return s
}

func (s *Scanner) Scan() (loc token.Location, kind token.Kind, text string) {
	s.skipWhitespace()

	if ch := s.ch; canLeadIdent(ch) {
		kind = token.KindIdent
		text = s.scanIdent()
		if len(text) > 1 {
			kind = token.Lookup(text)
		}
	} else if isDigit(ch) {
		kind, text = s.scanNumber()
	} else if ch == '"' {
		kind, text = token.KindString, s.scanString()
	} else if ch == '\'' {
		kind, text = token.KindChar, s.scanChar()
	} else if ch == '#' {
		kind, text = token.KindComment, s.scanComment()
	} else {
		// Operators do not need text
		switch ch {
		case eof:
			kind = token.KindEOF
		case '+':
			kind = token.KindAdd
		case '-':
			kind = token.KindSub
		case '*':
			kind = token.KindMul
		case '/':
			kind = token.KindDiv
			if s.expectNext('=') {
				kind = token.KindNe
			}
		case '%':
			kind = token.KindMod
		case '(':
			kind = token.KindLParen
		case ')':
			kind = token.KindRParen
		case '{':
			kind = token.KindLBrace
		case '}':
			kind = token.KindRBrace
		case '=':
			kind = token.KindAssign
			if s.expectNext('=') {
				kind = token.KindEq
			}
		case '>':
			kind = token.KindGt
			if s.expectNext('=') {
				kind = token.KindGe
			}
		case '<':
			kind = token.KindLt
			if s.expectNext('=') {
				kind = token.KindLe
			}
		case ';':
			kind = token.KindSemicolon
		case ',':
			kind = token.KindComma
		default:
			if ch != bom {
				s.reportf("illegal character %#U", ch)
			}
			kind, text = token.KindInvalid, string(ch)
		}
		s.next()
	}

	return
}

func (s *Scanner) expectNext(ch rune) bool {
	s.next()
	if s.ch == ch {
		s.next()
		return true
	}
	return false
}

func (s *Scanner) consume(ch rune) {
	if s.ch != ch {
		panic(fmt.Sprintf("expect current character: %q, actual: %q", ch, s.ch))
	}
	_ = s.next()
}

// take returns current character and advances
func (s *Scanner) take() (ch rune) {
	ch = s.ch
	_ = s.next()
	return ch
}

func (s *Scanner) scanIdent() string {
	off := s.offset
	for {
		ch := s.ch
		if !canMakeIdent(ch) {
			if canTerminateIdent(ch) {
				s.next()
			}
			break
		}
		s.next()
	}
	return string(s.src[off:s.offset])
}

func (s *Scanner) scanNumber() (kind token.Kind, text string) {
	kind = token.KindInt

	begin := s.offset

	if validFn := s.chooseAlphabet(); validFn != nil {
		// skip the first digit
		s.next()
		s.scanDigits(validFn)
	}

	if s.ch != '.' && s.ch != 'e' && s.ch != 'E' {
		return kind, string(s.src[begin:s.offset])
	}

	if ch := s.ch; ch == '.' {
		s.next()
		if n := s.scanDecimalDigits(); n == 0 {
			s.report("invalid floating-point literal: no fraction")
			return
		}
	}

	if ch := s.ch; ch == 'e' || ch == 'E' {
		s.next()
		if s.ch == '+' || s.ch == '-' {
			s.next()
		}
		if n := s.scanDecimalDigits(); n == 0 {
			s.report("invalid floating-point literal: incomplete exponent")
			return
		}
	}

	return token.KindFloat, string(s.src[begin:s.offset])
}

func (s *Scanner) chooseAlphabet() func(rune) bool {
	if s.ch != '0' {
		return isDigit
	}
	s.next()
	switch s.ch {
	case 'b':
		return isBinDigit
	case 'o':
		return isOctDigit
	case 'x':
		return isHexDigit
	}
	if isDigit(s.ch) {
		return isDigit
	}
	return nil
}

func (s *Scanner) scanDigits(isValid func(rune) bool) (n int) {
	for ch := s.ch; isValid(ch); ch = s.ch {
		n++
		s.next()
	}
	return
}

func (s *Scanner) scanDecimalDigits() (n int) {
	return s.scanDigits(isDigit)
}

// TODO: escape characters
func (s *Scanner) scanChar() string {
	begin := s.offset
	s.consume('\'')
	ch := s.take()
	if ch != '\'' {
		ch = s.take()
		if ch != '\'' {
			s.report("char literal not terminated")
		}
	} else {
		s.report("illegal char literal")
	}
	return string(s.src[begin:s.offset])
}

// TODO: escape characters
func (s *Scanner) scanString() string {
	begin := s.offset
	s.consume('"')
	for ch := s.ch; ch != '"'; ch = s.next() {
		if ch == '\n' || ch < 0 {
			s.report("string literal not terminated")
			break
		}
	}
	if s.ch == '"' {
		s.consume('"')
	}

	return string(s.src[begin:s.offset])
}

func (s *Scanner) scanComment() string {
	begin := s.offset
	s.consume('#')
	for ch := s.ch; ch != '\n' && ch > 0; ch = s.next() {
	}
	end := s.offset
	s.next()
	return string(s.src[begin:end])
}

func (s *Scanner) next() rune {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset

		// defaults to the case of ASCII
		r, n := rune(s.src[s.rdOffset]), 1
		if r == 0 {
			s.report("illegal character NUL")
		} else if r >= utf8.RuneSelf {
			// non-ASCII character
			r, n = utf8.DecodeRune(s.src[s.rdOffset:])
			if r == utf8.RuneError {
				switch n {
				case 0:
				case 1:
					s.report("illegal UTF-8 encoding")
				default:
					panic(fmt.Sprint("un-handled return value: ", n))
				}
			} else if r == bom && s.offset > 0 {
				s.report("illegal byte order mark")
			}
		}
		s.rdOffset += n
		s.ch = r
	} else {
		s.offset = len(s.src)
		s.ch = eof
	}

	return s.ch
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\r' || s.ch == '\n' {
		s.next()
	}
}

func canLeadIdent(ch rune) bool {
	return isLetter(ch) || ch == '_'
}

func canMakeIdent(ch rune) bool {
	return canLeadIdent(ch) || isDigit(ch)
}

var isDigit = isDecDigit

func isLetter(ch rune) bool {
	return ('A' <= ch && ch <= 'Z') || ('a' <= ch && ch <= 'z')
}

func canTerminateIdent(ch rune) bool {
	return ch == '!' || ch == '?'
}

func isBinDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func isOctDigit(ch rune) bool {
	return '0' <= ch && ch <= '7'
}

func isDecDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isHexDigit(ch rune) bool {
	return ('0' <= ch && ch <= '0') || ('a' <= ch && ch <= 'f') || ('A' <= ch && ch <= 'F')
}

func (s *Scanner) report(_ string) {
	s.NumErrors++
}

func (s *Scanner) reportf(format string, args ...interface{}) {
	s.report(fmt.Sprintf(format, args...))
}
