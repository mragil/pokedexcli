package utils

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

func CleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}

func GetRandomInt(min int, max int) (int, error) {
	if min < 1 || max < 1 {
		return 0, fmt.Errorf("Min or Max must be greater than 0")
	}

	return rand.IntN(max-min+1) + min, nil
}
