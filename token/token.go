package token

import "fmt"

var kindsByText = map[string]Kind{
	"true":  KindTrue,
	"false": KindFalse,
	"and":   KindAnd,
	"or":    KindOr,
	"not":   KindNot,

	"let": KindLet,

	"print": KindPrint,
}

func Lookup(ident string) Kind {
	if kind, ok := kindsByText[ident]; ok {
		return kind
	}
	return KindIdent
}

type Kind int

const (
	KindInvalid Kind = iota
	KindEOF
	KindComment

	literal_begin
	KindInt
	KindFloat
	KindChar
	KindString
	KindIdent
	literal_end

	operator_begin
	// Operators
	KindAdd // +
	KindSub // -
	KindMul // *
	KindDiv // /
	KindMod // %

	KindLParen // (
	KindRParen // )

	KindEq // ==
	KindNe // /=
	KindGt // >
	KindGe // >=
	KindLt // <
	KindLe // <=

	KindAssign // =

	KindSemicolon // ;
	KindComma     // ,
	operator_end

	keyword_begin
	// Keywords
	KindTrue  // true
	KindFalse // false
	KindAnd   // and
	KindOr    // or
	KindNot   // not

	KindLet // let

	KindPrint // print
	keyword_end
)

func (kind Kind) String() string {
	switch kind {
	case KindInvalid:
		return "<INVALID>"
	case KindEOF:
		return "EOF"
	case KindComment:
		return "COMMENT"

	case KindInt:
		return "INT"
	case KindFloat:
		return "FLOAT"
	case KindChar:
		return "CHAR"
	case KindString:
		return "STRING"
	case KindIdent:
		return "IDENT"

	case KindAdd:
		return "ADD"
	case KindSub:
		return "Sub"
	case KindMul:
		return "MUL"
	case KindDiv:
		return "Div"
	case KindMod:
		return "MOD"

	case KindLParen:
		return "LPAREN"
	case KindRParen:
		return "RPAREN"

	case KindEq:
		return "EQ"
	case KindNe:
		return "NE"
	case KindGt:
		return "GT"
	case KindGe:
		return "GE"
	case KindLt:
		return "LT"
	case KindLe:
		return "LE"

	case KindSemicolon:
		return "SEMICOLON"
	case KindComma:
		return "COMMA"

	case KindTrue:
		return "TRUE"
	case KindFalse:
		return "FALSE"
	case KindAnd:
		return "AND"
	case KindOr:
		return "OR"
	case KindNot:
		return "NOT"

	case KindPrint:
		return "PRINT"

	default:
		panic(fmt.Sprint("unknown token kind value: ", int(kind)))
	}
}

func init() {
	sz, nKeywordsDefined := len(kindsByText), int(keyword_end)-int(keyword_begin)-1
	if sz != nKeywordsDefined {
		panic(fmt.Sprintf("set keywords has %v elements but %v Kind's of keyword are defined", sz, nKeywordsDefined))
	}
}
