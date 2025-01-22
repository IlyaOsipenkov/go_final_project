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
	} else {
		cleanedRepeat := strings.ReplaceAll(repeat, ",", "")
		arrOfRuleAndDates := strings.Split(cleanedRepeat, " ")
		rule := arrOfRuleAndDates[0]
		daysOfRule := arrOfRuleAndDates[1:]

		switch rule {
		case "d":
			if len(arrOfRuleAndDates) < 2 {
				return "", errors.New("repeat doesnt valid for rule -d-")
			}
			days, err := strconv.Atoi(daysOfRule[0])
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
			daysOfWeek, err := parseDaysOfWeek(daysOfRule)
			if err != nil {
				return "", err
			}
			nextDate := getNextWeekday(now, daysOfWeek)
			return nextDate.Format("20060102"), nil

		case "m":

		default:
			fmt.Println("wrong rule")
		}
	}
	//написать сплит repeat и проверку на 2 часть
}

func parseDaysOfWeek(data []string) ([]int, error) {
	weekdays := []int{}

	for _, val := range data {
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
func main() {

}
