package main

import (
	"bytes"
	"testing"
)

func TestParseJson(t *testing.T) {
	testCases := []struct {
		data          []byte
		expectedError error
	}{
		{[]byte("  null "), PARSE_ROOT_NOT_SINGULAR},
		{[]byte("true"), nil},
		{[]byte("false"), nil},
		{[]byte("null"), nil},
	}
	for _, testCase := range testCases {
		var value jsonValue
		if err := parseJson(testCase.data, &value); err != testCase.expectedError {
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
		{[]byte("null"), jsonValue{NULL_TYPE}, []byte(""), nil},
		{[]byte("null "), jsonValue{NULL_TYPE}, []byte(" "), nil},
		{[]byte("nul"), jsonValue{0}, []byte("nul"), PARSE_INVALID_VALUE},
		{[]byte(""), jsonValue{0}, []byte(""), PARSE_INVALID_VALUE},
		{[]byte("nun"), jsonValue{0}, []byte("nun"), PARSE_INVALID_VALUE},
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
		{[]byte("true"), jsonValue{TRUE_TYPE}, []byte(""), nil},
		{[]byte("true "), jsonValue{TRUE_TYPE}, []byte(" "), nil},
		{[]byte("tru"), jsonValue{0}, []byte("tru"), PARSE_INVALID_VALUE},
		{[]byte(""), jsonValue{0}, []byte(""), PARSE_INVALID_VALUE},
		{[]byte("tue"), jsonValue{0}, []byte("tue"), PARSE_INVALID_VALUE},
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
		{[]byte("false"), jsonValue{FALSE_TYPE}, []byte(""), nil},
		{[]byte("false "), jsonValue{FALSE_TYPE}, []byte(" "), nil},
		{[]byte("fal"), jsonValue{0}, []byte("fal"), PARSE_INVALID_VALUE},
		{[]byte(""), jsonValue{0}, []byte(""), PARSE_INVALID_VALUE},
		{[]byte("fse"), jsonValue{0}, []byte("fse"), PARSE_INVALID_VALUE},
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
			jsonValue{TRUE_TYPE},
			[]byte(""),
			nil,
		},
		{
			[]byte("false"),
			"false",
			FALSE_TYPE,
			jsonValue{FALSE_TYPE},
			[]byte(""),
			nil,
		},
		{
			[]byte("null"),
			"null",
			NULL_TYPE,
			jsonValue{NULL_TYPE},
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
