package tools_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dotkafx/tools"
)

func TestSecondsToString(t *testing.T) {
	require := assert.New(t)

	testCases := map[string]struct {
		input          int
		requiredOutput string
	}{
		"nullValue": {
			0,
			"00:00:00",
		},
		"minusNullValue": {
			0,
			"00:00:00",
		},
		"oneValue": {
			1,
			"00:00:01",
		},
		"minusOneValue": {
			-1,
			"-00:00:01",
		},
		"hugeValue": {
			64000,
			"17:46:40",
		},
		"hugeNegativeValue": {
			-64000,
			"-17:46:40",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Logf("Testing SecondsToString, with %s", testCaseName)
		require.Equal(testCase.requiredOutput, tools.SecondsToString(testCase.input))
	}
}

func TestStringToSeconds(t *testing.T) {
	require := assert.New(t)

	testCases := map[string]struct {
		input          string
		requiredOutput int
		requiredError  string
	}{
		"invalidInput": {
			"invalid",
			0,
			`time: invalid duration "invalid"`,
		},
		"nullValue": {
			"0",
			0,
			"",
		},
		"minusNullValue": {
			"-0",
			0,
			"",
		},
		"oneValue": {
			"1",
			1,
			"",
		},
		"oneValueWithSeconds": {
			"1s",
			1,
			"",
		},
		"minusOneValue": {
			"-1",
			-1,
			"",
		},
		"minusOneValueWithSeconds": {
			"-1s",
			-1,
			"",
		},
		"hugeValue": {
			"17h46m40s",
			64000,
			"",
		},
		"hugeNegativeValue": {
			"-17h46m40s",
			-64000,
			"",
		},
		"nonClockValues": {
			"-25h61m61s",
			-93721,
			"",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Logf("Testing StringToSeconds, with %s", testCaseName)
		actualOutput, actualError := tools.StringToSeconds(testCase.input)
		require.Equal(testCase.requiredOutput, actualOutput)
		if testCase.requiredError == "" {
			require.NoError(actualError)
		} else {
			require.EqualError(actualError, testCase.requiredError)
		}
	}
}

func TestParseSuffixAmount(t *testing.T) {
	require := assert.New(t)

	testCases := map[string]struct {
		input          string
		inputPrefix    string
		requiredOutput int
		requiredError  string
	}{
		"invalidInput": {
			"invalid_invalid",
			"invalid",
			0,
			`time: invalid duration "_invalid"`,
		},
		"invalidInputNoPrefix": {
			"",
			"a",
			0,
			"Prefix a is missing.",
		},
		"nullValue": {
			"",
			"",
			1,
			"",
		},
		"minusNullValue": {
			"-0",
			"",
			0,
			"The amount value cannot be less than 1.",
		},
		"oneValue": {
			"a1",
			"a",
			1,
			"",
		},
		"twoValue": {
			"a2",
			"a",
			2,
			"",
		},
		"oneValueWithSeconds": {
			"a1s",
			"a",
			1,
			"",
		},
		"minusOneValue": {
			"a-1",
			"a",
			0,
			"The amount value cannot be less than 1.",
		},
		"minusTwoValueWithSeconds": {
			"a-2s",
			"a",
			0,
			"The amount value cannot be less than 1.",
		},
		"hugeValue": {
			"AA17h46m40s",
			"AA",
			0,
			"The mount value cannot be larger than 1800.",
		},
		"hugeNegativeValue": {
			"BB-17h46m40s",
			"BB",
			0,
			"The amount value cannot be less than 1.",
		},
		"nonClockValues": {
			"CC61s",
			"CC",
			61,
			"",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Logf("Testing StringToSeconds, with %s", testCaseName)
		actualOutput, actualError := tools.ParseSuffixAmount(testCase.input, testCase.inputPrefix)
		require.Equal(testCase.requiredOutput, actualOutput)
		if testCase.requiredError == "" {
			require.NoError(actualError)
		} else {
			require.EqualError(actualError, testCase.requiredError)
		}
	}
}
