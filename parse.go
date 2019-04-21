package main

const (
	_ = iota
	NULL_TYPE
	TRUE_TYPE
	FALSE_TYPE
	NUMBER_TYPE
	STRING_TYPE
	OBJECT_TYPE
	ARRAY_TYPE
)

// 表示 JSON 值
type jsonValue struct {
	_type int
}

// 解析 JSON
func parseJson(data []byte, value *jsonValue) error {
	data = parseWhitespace(data)
	var err error
	if bs, err := parseValue(data, value); err == nil {
		bs := parseWhitespace(bs)
		if len(bs) != 0 {
			err = PARSE_ROOT_NOT_SINGULAR
		}
	}
	return err
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
	return parseLiteral(data, "true", TRUE_TYPE, value)
}

func parseFalse(data []byte, value *jsonValue) ([]byte, error) {
	return parseLiteral(data, "false", FALSE_TYPE, value)
}

func parseNull(data []byte, value *jsonValue) ([]byte, error) {
	return parseLiteral(data, "null", NULL_TYPE, value)
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
	default:
		{
			return data, PARSE_INVALID_VALUE
		}
	}
}
