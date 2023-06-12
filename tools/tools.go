package tools

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseSuffixAmount is used to parse a string like "back123" by trimming the prefix and converting the suffix int integer
func ParseSuffixAmount(text string, prefix string) (int, error) {
	if !strings.HasPrefix(text, prefix) {
		return 0, fmt.Errorf("Prefix %s is missing.", prefix)
	}
	if text == prefix {
		return 1, nil
	}
	amount, err := StringToSeconds(strings.TrimPrefix(text, prefix))
	if err != nil {
		return 0, err
	}
	if amount < 1 {
		return 0, fmt.Errorf("The amount value cannot be less than 1.")
	}
	if amount > 1800 {
		return 0, fmt.Errorf("The mount value cannot be larger than 1800.")
	}
	return amount, nil
}

// SecondsToString returns a string in the hh:mm:ss format (e.g.: input: 4374 output: "01:12:54")
func SecondsToString(seconds int) string {
	prefix := ""
	if seconds < 0 {
		prefix = "-"
	}

	hours := absInt(seconds / 3600 % 24)
	minutes := absInt(seconds / 60 % 60)
	seconds = absInt(seconds % 60)

	return fmt.Sprintf("%s%02d:%02d:%02d", prefix, hours, minutes, seconds)
}

func absInt(x int) int {
	if x >= 0 {
		return x
	}

	return -x
}

func StringToSeconds(input string) (int, error) {
	if _, err := strconv.Atoi(input); err == nil {
		input = input + "s"
	}

	val, err := time.ParseDuration(input)
	if err != nil {
		return 0, err
	}
	return int(val.Seconds()), nil
}
