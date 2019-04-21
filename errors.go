package main

import "errors"

var (
	// 若 JSON 只含有空白字符，则返回此错误
	PARSE_EXPECT_VALUE = errors.New("expect value")
	// 若无法解析出合法的 JSON 值，则返回此错误
	PARSE_INVALID_VALUE = errors.New("invalid value")
	// 成功从 JSON 中解析出一个值，然后在空白字符之后还有其它字符，则返回此错误
	PARSE_ROOT_NOT_SINGULAR = errors.New("root not sigular")
)
