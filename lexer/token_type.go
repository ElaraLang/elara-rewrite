package lexer

type TokenType int

const (
	//special tokens
	Illegal TokenType = iota
	EOF
	NEWLINE

	//Brackets
	LParen
	RParen
	LBrace
	RBrace
	LAngle //<
	RAngle //>
	LSquare
	RSquare

	//Keywords
	Let
	Extend
	Return
	While
	Mut
	Lazy
	Restricted
	Struct
	Namespace
	Import
	Type
	If
	Else
	Match
	As
	Is

	//Operators
	Add
	Subtract
	Multiply
	Slash
	Mod
	And
	Or
	Xor
	Equals
	NotEquals
	GreaterEqual
	LesserEqual
	Not

	TypeOr  // |
	TypeAnd // &

	//Symbol
	Equal
	Arrow
	Dot

	//Literals
	BooleanTrue
	BooleanFalse
	String
	Char
	Int
	Float

	Comma
	Colon

	Identifier
	Underscore
)

func (token *TokenType) String() string {
	return tokenNames[*token]
}

var tokenNames = map[TokenType]string{
	Illegal: "Illegal",
	EOF:     "EOF",
	NEWLINE: "\\n",

	LParen:       "LParen",
	RParen:       "RParen",
	LBrace:       "LBrace",
	RBrace:       "RBrace",
	LAngle:       "LAngle",
	RAngle:       "RAngle",
	LSquare:      "LSquare",
	RSquare:      "RSquare",
	Type:         "Type",
	Let:          "Let",
	Extend:       "Extend",
	Return:       "Return",
	While:        "While",
	Mut:          "Mut",
	Lazy:         "Lazy",
	Restricted:   "Restricted",
	Struct:       "Struct",
	Namespace:    "Namespace",
	Import:       "Import",
	If:           "If",
	Else:         "Else",
	Match:        "Match",
	As:           "As",
	Is:           "Is",
	Add:          "Add",
	Subtract:     "Subtract",
	Multiply:     "Multiply",
	Slash:        "Slash",
	Mod:          "Mod",
	And:          "And",
	Or:           "Or",
	Xor:          "Xor",
	Equals:       "Equals",
	NotEquals:    "NotEquals",
	GreaterEqual: "GreaterEqual",
	LesserEqual:  "LesserEqual",
	Not:          "Not",
	Equal:        "Equal",
	Arrow:        "Arrow",
	Dot:          "Dot",
	BooleanTrue:  "True",
	BooleanFalse: "False",
	String:       "String",
	Char:         "Char",
	Int:          "Int",
	Float:        "Float",

	Comma: "Comma",
	Colon: "Colon",

	Identifier: "Identifier",
	Underscore: "Underscore",
}
var IllegalIdentifierChars = []bool{
	',':  true,
	'.':  true,
	':':  true,
	'#':  true,
	'[':  true,
	']':  true,
	'(':  true,
	')':  true,
	'{':  true,
	'}':  true,
	'"':  true,
	'>':  true,
	'<':  true,
	' ':  true,
	'\n': true,
	'\r': true,
	'\t': true,
}
