package utils

import (
	"fmt"
	"strconv"
)

const millisecondsInSecond = 1000

func SecondsToMilliseconds(seconds int64) int64 {
	return seconds * millisecondsInSecond
}

func ParseTimeParam(param string) (int64, error) {
	value, err := strconv.ParseInt(param, 10, 64)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid time param: %s", param)
	}
	return value, nil
}
