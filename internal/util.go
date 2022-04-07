package internal

import "fmt"

var RussianMonths = map[int]string{
	1:  "Январь",
	2:  "Февраль",
	3:  "Март",
	4:  "Апрель",
	5:  "Май",
	6:  "Июнь",
	7:  "Июль",
	8:  "Август",
	9:  "Сентябрь",
	10: "Октябрь",
	11: "Ноябрь",
	12: "Декабрь",
}

func GetRussianMonth(month int) (string, error) {
	if month < 1 || month > 12 {
		return "", fmt.Errorf("wrong month")
	}

	str := RussianMonths[month]

	return str, nil
}

func GetNextRussianMonth(currentMonth int) (string, error) {
	if currentMonth < 1 || currentMonth > 12 {
		return "", fmt.Errorf("wrong month")
	}

	var month string

	if currentMonth == 12 {
		month = RussianMonths[1]
	} else {
		month = RussianMonths[currentMonth+1]
	}

	return month, nil
}

func GetPastRussianMonth(currentMonth int) (string, error) {
	if currentMonth < 1 || currentMonth > 12 {
		return "", fmt.Errorf("wrong month")
	}

	var month string

	if currentMonth == 1 {
		month = RussianMonths[12]
	} else {
		month = RussianMonths[currentMonth-1]
	}

	return month, nil
}
