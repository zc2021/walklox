// Code generated by "stringer -type=TokenType ../tokens"; DO NOT EDIT.

package tokens

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[EOF-0]
	_ = x[begin_delimiters-1]
	_ = x[LEFT_PAREN-2]
	_ = x[RIGHT_PAREN-3]
	_ = x[LEFT_BRACE-4]
	_ = x[RIGHT_BRACE-5]
	_ = x[COMMA-6]
	_ = x[DOT-7]
	_ = x[SEMICOLON-8]
	_ = x[end_delimiters-9]
	_ = x[begin_literals-10]
	_ = x[IDENTIFIER-11]
	_ = x[STRING-12]
	_ = x[NUMBER-13]
	_ = x[end_literals-14]
	_ = x[begin_binary_operators-15]
	_ = x[PLUS-16]
	_ = x[SLASH-17]
	_ = x[STAR-18]
	_ = x[EQUAL_EQUAL-19]
	_ = x[EQUAL-20]
	_ = x[GREATER_EQUAL-21]
	_ = x[GREATER-22]
	_ = x[LESS_EQUAL-23]
	_ = x[LESS-24]
	_ = x[BANG_EQUAL-25]
	_ = x[end_binary_operators-26]
	_ = x[begin_unary_operators-27]
	_ = x[MINUS-28]
	_ = x[end_undary_operators-29]
	_ = x[begin_variadic_operators-30]
	_ = x[BANG-31]
	_ = x[end_variadic_operators-32]
	_ = x[begin_reserved_keywords-33]
	_ = x[AND-34]
	_ = x[CLASS-35]
	_ = x[ELSE-36]
	_ = x[FALSE-37]
	_ = x[FUN-38]
	_ = x[FOR-39]
	_ = x[IF-40]
	_ = x[NIL-41]
	_ = x[OR-42]
	_ = x[PRINT-43]
	_ = x[RETURN-44]
	_ = x[SUPER-45]
	_ = x[THIS-46]
	_ = x[TRUE-47]
	_ = x[VAR-48]
	_ = x[WHILE-49]
	_ = x[end_reserved_keywords-50]
}

const _TokenType_name = "EOFbegin_delimitersLEFT_PARENRIGHT_PARENLEFT_BRACERIGHT_BRACECOMMADOTSEMICOLONend_delimitersbegin_literalsIDENTIFIERSTRINGNUMBERend_literalsbegin_binary_operatorsPLUSSLASHSTAREQUAL_EQUALEQUALGREATER_EQUALGREATERLESS_EQUALLESSBANG_EQUALend_binary_operatorsbegin_unary_operatorsMINUSend_undary_operatorsbegin_variadic_operatorsBANGend_variadic_operatorsbegin_reserved_keywordsANDCLASSELSEFALSEFUNFORIFNILORPRINTRETURNSUPERTHISTRUEVARWHILEend_reserved_keywords"

var _TokenType_index = [...]uint16{0, 3, 19, 29, 40, 50, 61, 66, 69, 78, 92, 106, 116, 122, 128, 140, 162, 166, 171, 175, 186, 191, 204, 211, 221, 225, 235, 255, 276, 281, 301, 325, 329, 351, 374, 377, 382, 386, 391, 394, 397, 399, 402, 404, 409, 415, 420, 424, 428, 431, 436, 457}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
