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

// ParseSet parses:
//
//	/set TP 20
//
// returning:
//
//	systemID = "TP"
//	amount = 20
func ParseSet(args string) (systemID string, amount int, err error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", 0, fmt.Errorf("usage: /set <SYSTEM> <AMOUNT>")
	}

	amount, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("amount must be an integer")
	}

	if amount < 0 {
		return "", 0, fmt.Errorf("amount must be non-negative")
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

// ParseRegen parses:
//
//	/regen TP 15
//
// returning:
//
//	systemID = "TP"
//	minutesLeft = 15
func ParseRegen(args string) (systemID string, minutesLeft int, err error) {
	fields := strings.Fields(args)
	if len(fields) != 2 {
		return "", 0, fmt.Errorf("usage: /regen <SYSTEM> <MINUTES_LEFT>")
	}

	minutesLeft, err = strconv.Atoi(fields[1])
	if err != nil {
		return "", 0, fmt.Errorf("minutes must be an integer")
	}

	if minutesLeft < 0 {
		return "", 0, fmt.Errorf("minutes must be non-negative")
	}

	return strings.ToUpper(fields[0]), minutesLeft, nil
}
