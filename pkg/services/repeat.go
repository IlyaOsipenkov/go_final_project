package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	parsedDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid parse of date: %v", err)
	}

	if len(repeat) == 0 {
		fmt.Println("delete task from db")
		return "", fmt.Errorf("repeat is empty. Length of repeat must be greater than 0")
	}

	arrOfRuleAndDates := strings.Fields(repeat)
	rule := arrOfRuleAndDates[0]

	switch rule {
	case "d":
		if len(arrOfRuleAndDates) != 2 {
			return "", errors.New("repeat doesnt valid for rule -d-")
		}
		days, err := strconv.Atoi(arrOfRuleAndDates[1])
		if err != nil || days < 0 || days > 400 {
			return "", errors.New("days of repead doesnt valid. Must be in 1-400")
		}
		nextDate := parsedDate.AddDate(0, 0, days)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil
	case "y":
		if len(arrOfRuleAndDates) > 1 {
			return "", errors.New("repeat doesnt valid for rule -y-")
		}
		nextDate := parsedDate.AddDate(1, 0, 0)
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil
	case "w":
		if len(arrOfRuleAndDates) != 2 {
			return "", errors.New("repeat doesnt valid for rule -w-")
		}
		daysOfWeek, err := parseDaysOfWeek(arrOfRuleAndDates[1])
		if err != nil {
			return "", err
		}
		nextDate := getNextWeekday(parsedDate, daysOfWeek)
		for nextDate.Before(now) {
			nextDate = getNextWeekday(nextDate.AddDate(0, 0, 1), daysOfWeek)
		}
		return nextDate.Format("20060102"), nil
	default:
		return "", errors.New("invalid rule")
	}
}

func parseDaysOfWeek(data string) ([]int, error) {
	days := strings.Split(data, ",")
	weekdays := []int{}
	for _, val := range days {
		day, err := strconv.Atoi(val)
		if err != nil || day < 1 || day > 7 {
			return nil, fmt.Errorf("the day doesnt valid: %s", val)
		}
		weekdays = append(weekdays, day)
	}
	return weekdays, nil

}

func getNextWeekday(start time.Time, weekdays []int) time.Time {
	for i := 0; i < 7; i++ {
		dayOfWeekNow := int(start.Weekday())
		if dayOfWeekNow == 0 {
			dayOfWeekNow = 7
		}
		for _, weekday := range weekdays {
			if dayOfWeekNow == weekday && start.After(time.Now()) {
				return start
			}
		}
		start = start.AddDate(0, 0, 1)
	}
	return start
}
