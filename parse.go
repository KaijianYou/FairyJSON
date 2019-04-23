package main

import (
	"strconv"
)

const (
	INVALID_TYPE = iota
	NULL_TYPE
	TRUE_TYPE
	FALSE_TYPE
	NUMBER_TYPE
	STRING_TYPE
	OBJECT_TYPE
	ARRAY_TYPE
)

const (
	NUMBER_LITERAL = "0123456789"
	TRUE_LITERAL   = "true"
	FALSE_LITERAL  = "false"
	NULL_LITERAL   = "null"
)

// 表示 JSON 值
type jsonValue struct {
	number float64 // 仅当 _type 为 NUMBER_TYPE 时，number 才表示 JSON 数值类型的数值
	_type  int
}

// 解析 JSON
func parseJson(data []byte) (jsonValue, error) {
	var value jsonValue
	var err error
	data = parseWhitespace(data)
	bs, err := parseValue(data, &value)
	if err == nil {
		bs := parseWhitespace(bs)
		if len(bs) != 0 {
			err = PARSE_ROOT_NOT_SINGULAR
		}
	}
	return value, err
}

// 解析（忽略）空白符，只允许出现空格符、制表符、换行符和回车符
func parseWhitespace(data []byte) []byte {
	i := 0
	for ; i < len(data); i++ {
		if !(data[i] == ' ' || data[i] == '\t' || data[i] == '\n' || data[i] == '\r') {
			break
		}
	}
	return data[i:]
}

func parseLiteral(data []byte, literal string, jsonType int, value *jsonValue) ([]byte, error) {
	if len(data) < len(literal) {
		return data, PARSE_INVALID_VALUE
	}
	if string(data[:len(literal)]) != literal {
		return data, PARSE_INVALID_VALUE
	}
	value._type = jsonType
	return data[len(literal):], nil
}

func parseTrue(data []byte, value *jsonValue) ([]byte, error) {
	return parseLiteral(data, TRUE_LITERAL, TRUE_TYPE, value)
}

func parseFalse(data []byte, value *jsonValue) ([]byte, error) {
	return parseLiteral(data, FALSE_LITERAL, FALSE_TYPE, value)
}

func parseNull(data []byte, value *jsonValue) ([]byte, error) {
	return parseLiteral(data, NULL_LITERAL, NULL_TYPE, value)
}

func getNumber(value jsonValue) (float64, error) {
	if value._type != NUMBER_TYPE {
		return 0, REQUIRE_NUMBER_TYPE
	}
	return value.number, nil
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func isDigit1To9(b byte) bool {
	return b >= '1' && b <= '9'
}

func parseNumber(data []byte, value *jsonValue) ([]byte, error) {
	i := 0
	bs := append(data[:], '\x00')
	if bs[i] == '-' {
		i++
	}
	if bs[i] == '0' {
		i++
	} else {
		if !isDigit1To9(bs[i]) {
			return bs, PARSE_INVALID_VALUE
		}
		for isDigit(bs[i]) {
			i++
		}
	}
	if bs[i] == '.' {
		i++
		if !isDigit(bs[i]) {
			return bs, PARSE_INVALID_VALUE
		}
		for isDigit(bs[i]) {
			i++
		}
	}
	if bs[i] == 'e' || bs[i] == 'E' {
		i++
		if bs[i] == '+' || bs[i] == '-' {
			i++
		}
		if !isDigit(bs[i]) {
			return bs, PARSE_INVALID_VALUE
		}
		for isDigit(bs[i]) {
			i++
		}
	}
	var err error
	value.number, err = strconv.ParseFloat(string(data), 64)
	if err != nil {
		return data, PARSE_INVALID_VALUE
	}
	value._type = NUMBER_TYPE
	return data, nil
}

func parseValue(data []byte, value *jsonValue) ([]byte, error) {
	if data == nil || len(data) == 0 {
		return []byte{}, PARSE_EXPECT_VALUE
	}
	switch data[0] {
	case 'n':
		{
			return parseNull(data, value)
		}
	case 't':
		{
			return parseTrue(data, value)
		}
	case 'f':
		{
			return parseFalse(data, value)
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
		{
			return parseNumber(data, value)
		}
	default:
		{
			return data, PARSE_INVALID_VALUE
		}
	}
}
