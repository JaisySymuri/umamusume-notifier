package telegram

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseUse parses:
//
//	/use TP 20
//
// returning:
//
//	systemID = "TP"
//	amount = 20
func ParseUse(args string) (systemID string, amount int, err error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", 0, fmt.Errorf("usage: /use <SYSTEM> <AMOUNT>")
	}

	amount, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("amount must be an integer")
	}

	return strings.ToUpper(fields[0]), amount, nil
}

// ParseElapsed parses:
//
//	/elapsed TP 15
//
// returning:
//
//	systemID = "TP"
//	minutes = 15
func ParseElapsed(args string) (systemID string, minutes int, err error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", 0, fmt.Errorf("usage: /elapsed <SYSTEM> <MINUTES>")
	}

	minutes, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("minutes must be an integer")
	}

	return strings.ToUpper(fields[0]), minutes, nil
}