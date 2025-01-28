package repeat

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
		return "", fmt.Errorf("Invalid parse of date: %v", err)
	}

	if len(repeat) == 0 {
		fmt.Println("delete task from db")
		return "", fmt.Errorf("Repeat is empty. Length of repeat must be greater than 0")
	} else {

		// cleanedRepeat := strings.ReplaceAll(repeat, ",", "")
		// arrOfRuleAndDates := strings.Split(cleanedRepeat, " ")
		// rule := arrOfRuleAndDates[0]
		// daysOfRule := arrOfRuleAndDates[1:]
		arrOfRuleAndDates := strings.Fields(repeat)
		rule := arrOfRuleAndDates[0]

		switch rule {
		case "d":
			if len(arrOfRuleAndDates) != 2 {
				return "", errors.New("repeat doesnt valid for rule -d-")
			}
			days, err := strconv.Atoi(arrOfRuleAndDates[1])
			if err != nil || days < 0 || days < 400 {
				return "", errors.New("days of repead doesnt valid. Must be in 1-400")
			}
			nextDate := parsedDate.AddDate(0, 0, days)
			if nextDate.Before(now) {
				return NextDate(now, nextDate.Format("20060102"), repeat)
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
		case "w":
			if len(arrOfRuleAndDates) != 2 {
				return "", errors.New("repeat doesnt valid for rule -w-")
			}
			daysOfWeek, err := parseDaysOfWeek(arrOfRuleAndDates[1])
			if err != nil {
				return "", err
			}
			nextDate := getNextWeekday(now, daysOfWeek)
			return nextDate.Format("20060102"), nil

		case "m":
			if len(arrOfRuleAndDates) < 2 {
				return "", fmt.Errorf("repeat doesnt valid for rule 'm'")
			}
			months, err := parseMonthesOfRule(arrOfRuleAndDates[1:])
			if err != nil {
				return "", err
			}
			nextDate = getNextMonthDate(now, months)
			return nextDate.Format("20060102"), nil
			//3. Проверить после удаления заглушки
		default:
			fmt.Println("wrong rule")
		}
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

func getNextWeekday(now time.Time, weekdays []int) time.Time {
	for {
		dayOfWeekNow := int(now.Weekday())
		if dayOfWeekNow == 0 {
			dayOfWeekNow = 7
		}
		for _, weekday := range weekdays {
			if dayOfWeekNow == weekday {
				return now
			}
		}
		now = now.AddDate(0, 0, 1)
	}
}

func parseMonthesOfRule(data []string) ([]int, error) {
	days := strings.Split(data[0], ",")
	//1. Дописать вызов функции для обработки дней days
	parsedMonths := []int{}
	if len(data) > 1 {
		months := strings.Split(data[1], ",")
		parsedMonths = parseStringMonthToInt(months)
		for _, m := range parsedMonths {
			if m < 1 || m > 12 {
				return nil, fmt.Errorf("month must be > 1 or < 12, but now its: %d", m)
			}
		}
	}
	return parsedMonths, nil
}

func parseStringMonthToInt(data []string) []int {
	months := []int{}
	for _, val := range data {
		v, err := strconv.Atoi(val)
		if err == nil {
			months = append(months, v)
		}

	}
	return months
}

func getNextMonthDate(now time.Time, days []int) time.Time {
	return now
	//Дописать заглушку
}
func main() {

}
