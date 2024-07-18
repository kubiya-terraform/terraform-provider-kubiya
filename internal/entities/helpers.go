package entities

import (
	"fmt"
	"strconv"
	"strings"
)

func isDaily(cronExpr string) bool {
	_, _, day, month, weekday, err := parseCron(cronExpr)
	if err != nil {
		return false
	}
	return day == "*" && month == "*" && weekday == "*"
}

func isWeekly(cronExpr string) bool {
	_, _, day, month, weekday, err := parseCron(cronExpr)
	if err != nil {
		return false
	}
	dayInt, err := strconv.Atoi(day)
	if err == nil && dayInt >= 0 && dayInt <= 6 {
		return false
	}
	return day == "*" && month == "*" && weekday != "*"
}

func isHourly(cronExpr string) bool {
	_, hour, day, month, weekday, err := parseCron(cronExpr)
	if err != nil {
		return false
	}
	return hour == "*" && day == "*" && month == "*" && weekday == "*"
}

func isMonthly(cronExpr string) bool {
	_, _, day, month, weekday, err := parseCron(cronExpr)
	if err != nil {
		return false
	}
	return day != "*" && month == "*" && weekday == "*"
}

func validCron(cronExpr string) error {
	_, _, _, _, _, err := parseCron(cronExpr)
	return err
}

func parseCron(cronExpr string) (minute, hour, day, month, weekday string, err error) {
	parts := strings.Fields(cronExpr)
	if len(parts) != 6 {
		return "", "", "", "", "", fmt.Errorf("invalid cron expression")
	}
	return parts[0], parts[1], parts[2], parts[3], parts[4], nil
}
