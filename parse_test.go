package main

import (
	"bytes"
	"math"
	"testing"
)

func TestParseJson(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedError error
	}{
		{[]byte("  null }"), PARSE_ROOT_NOT_SINGULAR},
		{[]byte("  null "), nil},
		{[]byte("true"), nil},
		{[]byte("false"), nil},
		{[]byte("null"), nil},
	}
	for _, testCase := range testCases {
		if _, err := parseJson(testCase.data); err != testCase.expectedError {
			t.Errorf(
				"Failed (TestParseJson) actual error: [%v], expected error: [%v]",
				err,
				testCase.expectedError,
			)
		}
	}
}

func TestParseValue(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedData  []byte
		expectedError error
	}{
		{[]byte("true"), []byte(""), nil},
		{[]byte("tr"), []byte("tr"), PARSE_INVALID_VALUE},
		{[]byte("false"), []byte(""), nil},
		{[]byte("fal"), []byte("fal"), nil},
		{[]byte("null"), []byte(""), nil},
		{[]byte("nu"), []byte("nu"), nil},
		{[]byte(""), []byte(""), PARSE_EXPECT_VALUE},
		{[]byte("o"), []byte("o"), PARSE_INVALID_VALUE},
	}
	for _, testCase := range testCases {
		var value jsonValue
		result, err := parseValue(testCase.data, &value)
		if !bytes.Equal(result, testCase.expectedData) {
			t.Errorf(
				"Failed (TestParseValue) actual data: [%v], expected data: [%v]\tactual error: [%v], expected error: [%v]",
				result,
				testCase.expectedData,
				err,
				testCase.expectedError,
			)
		}
	}
}

func TestParseWhitespace(t *testing.T) {
	testCases := []struct {
		data     []byte
		expected []byte
	}{
		{[]byte(" \t \t\r \n abc \t "), []byte("abc \t ")},
		{[]byte(" \t \t\r \n abc"), []byte("abc")},
		{[]byte(""), []byte("")},
	}
	for _, testCase := range testCases {
		result := parseWhitespace(testCase.data)
		if !bytes.Equal(result, testCase.expected) {
			t.Errorf(
				"Failed (TestParseWhitespace) actual: [%v], expected: [%v]",
				result,
				testCase.expected,
			)
		}
	}
}

func TestParseNull(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedValue jsonValue
		expectedData  []byte
		expectedError error
	}{
		{[]byte("null"), jsonValue{0, NULL_TYPE}, []byte(""), nil},
		{[]byte("null "), jsonValue{0, NULL_TYPE}, []byte(" "), nil},
		{[]byte("nul"), jsonValue{0, INVALID_TYPE}, []byte("nul"), PARSE_INVALID_VALUE},
		{[]byte(""), jsonValue{0, INVALID_TYPE}, []byte(""), PARSE_INVALID_VALUE},
		{[]byte("nun"), jsonValue{0, INVALID_TYPE}, []byte("nun"), PARSE_INVALID_VALUE},
	}
	for _, testCase := range testCases {
		value := jsonValue{}
		bs, err := parseNull(testCase.data, &value)
		if !bytes.Equal(bs, testCase.expectedData) || value != testCase.expectedValue || err != testCase.expectedError {
			t.Errorf(
				"Failed (TestParseNull) actual data: [%v], expected data: [%v]\tactual value: [%v], expected value: [%v]\tactual error: [%v], expected error: [%v]",
				bs,
				testCase.expectedData,
				value,
				testCase.expectedValue,
				err,
				testCase.expectedError,
			)
		}
	}
}

func TestParseTrue(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedValue jsonValue
		expectedData  []byte
		expectedError error
	}{
		{[]byte("true"), jsonValue{0, TRUE_TYPE}, []byte(""), nil},
		{[]byte("true "), jsonValue{0, TRUE_TYPE}, []byte(" "), nil},
		{[]byte("tru"), jsonValue{0, INVALID_TYPE}, []byte("tru"), PARSE_INVALID_VALUE},
		{[]byte(""), jsonValue{0, INVALID_TYPE}, []byte(""), PARSE_INVALID_VALUE},
		{[]byte("tue"), jsonValue{0, INVALID_TYPE}, []byte("tue"), PARSE_INVALID_VALUE},
	}
	for _, testCase := range testCases {
		value := jsonValue{}
		bs, err := parseTrue(testCase.data, &value)
		if !bytes.Equal(bs, testCase.expectedData) || value != testCase.expectedValue || err != testCase.expectedError {
			t.Errorf(
				"Failed (TestParseTrue) actual data: [%v], expected data: [%v]\tactual value: [%v], expected value: [%v]\tactual error: [%v], expected error: [%v]",
				bs,
				testCase.expectedData,
				value,
				testCase.expectedValue,
				err,
				testCase.expectedError,
			)
		}
	}
}

func TestParseFalse(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedValue jsonValue
		expectedData  []byte
		expectedError error
	}{
		{[]byte("false"), jsonValue{0, FALSE_TYPE}, []byte(""), nil},
		{[]byte("false "), jsonValue{0, FALSE_TYPE}, []byte(" "), nil},
		{[]byte("fal"), jsonValue{0, INVALID_TYPE}, []byte("fal"), PARSE_INVALID_VALUE},
		{[]byte(""), jsonValue{0, INVALID_TYPE}, []byte(""), PARSE_INVALID_VALUE},
		{[]byte("fse"), jsonValue{0, INVALID_TYPE}, []byte("fse"), PARSE_INVALID_VALUE},
	}
	for _, testCase := range testCases {
		value := jsonValue{}
		bs, err := parseFalse(testCase.data, &value)
		if !bytes.Equal(bs, testCase.expectedData) || value != testCase.expectedValue || err != testCase.expectedError {
			t.Errorf(
				"Failed (TestParseFalse) actual data: [%v], expected data: [%v]\tactual value: [%v], expected value: [%v]\tactual error: [%v], expected error: [%v]",
				bs,
				testCase.expectedData,
				value,
				testCase.expectedValue,
				err,
				testCase.expectedError,
			)
		}
	}
}

func TestParseLiteral(t *testing.T) {
	testCases := []struct {
		data          []byte
		literal       string
		jsonType      int
		expectedValue jsonValue
		expectedData  []byte
		expectedError error
	}{
		{
			[]byte("true"),
			"true",
			TRUE_TYPE,
			jsonValue{0, TRUE_TYPE},
			[]byte(""),
			nil,
		},
		{
			[]byte("false"),
			"false",
			FALSE_TYPE,
			jsonValue{0, FALSE_TYPE},
			[]byte(""),
			nil,
		},
		{
			[]byte("null"),
			"null",
			NULL_TYPE,
			jsonValue{0, NULL_TYPE},
			[]byte(""),
			nil,
		},
	}
	for _, testCase := range testCases {
		value := jsonValue{}
		bs, err := parseLiteral(testCase.data, testCase.literal, testCase.jsonType, &value)
		if !bytes.Equal(bs, testCase.expectedData) || value != testCase.expectedValue || err != testCase.expectedError {
			t.Errorf(
				"Failed (TestParseLiteral) actual data: [%v], expected data: [%v]\tactual value: [%v], "+
					"expected value: [%v]\tactual error: [%v], expected error: [%v]",
				bs,
				testCase.expectedData,
				value,
				testCase.expectedValue,
				err,
				testCase.expectedError,
			)
		}
	}
}

func TestParseNumber(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedValue jsonValue
		expectedError error
	}{
		{
			[]byte("0"),
			jsonValue{0.0, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-0"),
			jsonValue{0.0, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-0.0"),
			jsonValue{0.0, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1"),
			jsonValue{1.0, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-1"),
			jsonValue{-1.0, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1.5"),
			jsonValue{1.5, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-1.5"),
			jsonValue{-1.5, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("3.1415926"),
			jsonValue{3.1415926, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1E10"),
			jsonValue{1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1e10"),
			jsonValue{1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-1E10"),
			jsonValue{-1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-1e10"),
			jsonValue{-1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-1E+10"),
			jsonValue{-1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1E+10"),
			jsonValue{1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1e+10"),
			jsonValue{1e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("-1e-10"),
			jsonValue{-1e-10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1.23E+10"),
			jsonValue{1.23e+10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1.23E-10"),
			jsonValue{1.23e-10, NUMBER_TYPE},
			nil,
		},
		{
			[]byte("1e10000"), // 溢出
			jsonValue{math.Inf(0), INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("+0"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("+1"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte(".123"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("1."),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("INF"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("inf"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("NAN"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("nan"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("1-1"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("-inf"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte(""),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		{
			[]byte("-.2"),
			jsonValue{0.0, INVALID_TYPE},
			PARSE_INVALID_VALUE,
		},
		// {
		// 	[]byte("1e-100000"), // 溢出
		// 	jsonValue{0.0, INVALID_TYPE},
		// 	PARSE_INVALID_VALUE,
		// },
	}
	for _, testCase := range testCases {
		value := jsonValue{}
		_, err := parseNumber(testCase.data, &value)
		if value != testCase.expectedValue || err != testCase.expectedError {
			t.Errorf(
				"Failed (TestParseNumber) actual value: [%v], expected value: [%v]\tactual error: [%v], expected error: [%v]",
				value,
				testCase.expectedValue,
				err,
				testCase.expectedError,
			)
		}
	}
}
