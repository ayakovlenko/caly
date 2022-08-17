package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

var monthYear = flag.String("t", "", "YYYY/MM")

func main() {
	flag.Parse()

	dow := []string{
		textBold("Mo"),
		textBold("Tu"),
		textBold("We"),
		textBold("Th"),
		textBold("Fr"),
		textBold("Sa"),
		textBold("Su"),
	}

	t, err := parseTime(*monthYear)
	if err != nil {
		panic(err)
	}

	cal := NewCal(t)

	// print month
	fmt.Printf(textBold("%s %d\n\n"), t.Month(), t.Year())

	// print days of week
	fmt.Println("\t" + strings.Join(dow, "\t"))

	// print days
	for _, w := range cal.weeks {
		s := []string{
			textBold(fmt.Sprintf("W%d", w[0])),
		}
		for _, d := range w[1:] {
			if d == 0 {
				s = append(s, "")
			} else {
				// highlight today's date
				if isToday(t, d) {
					s = append(s, textHighlight(fmt.Sprintf("%2d", d)))
					continue
				}

				s = append(s, fmt.Sprintf("%2d", d))
			}
		}

		fmt.Println(strings.Join(s, "\t"))
	}
	fmt.Println()
}

type cal struct {
	// [0] week number
	// [1:] dates
	weeks [][8]int
}

func NewCal(t time.Time) cal {
	year, month := t.Year(), t.Month()
	date := time.Date(year, month, 1, 0, 0, 0, 0, t.Local().Location())

	weeks := [][8]int{}

	for date.Month() == month {
		week := [8]int{}

		weekNum := getWeek(date)
		week[0] = weekNum

		for weekNum == getWeek(date) && month == date.Month() {
			week[getDay(date)] = date.Day()
			date = date.AddDate(0, 0, 1)
		}
		weeks = append(weeks, week)
	}

	return cal{weeks}
}

func getWeek(t time.Time) int {
	_, w := t.ISOWeek()
	return w
}

var weekdayMap map[time.Weekday]int = map[time.Weekday]int{
	time.Monday:    1,
	time.Tuesday:   2,
	time.Wednesday: 3,
	time.Thursday:  4,
	time.Friday:    5,
	time.Saturday:  6,
	time.Sunday:    7,
}

func getDay(t time.Time) int {
	return weekdayMap[t.Weekday()]
}

// constants to control textDecorate
const (
	FontBold    = "\u001b[1m"
	FontReverse = "\u001b[7m"
)

func textDecorate(mod string, s string) string {
	return mod + s + "\u001b[0m"
}

func textBold(s string) string {
	return textDecorate(FontBold, s)
}

func textHighlight(s string) string {
	return textDecorate(FontReverse, s)
}

func isToday(t time.Time, day int) bool {
	now := time.Now()
	var (
		sameYear  = now.Year() == t.Year()
		sameMonth = now.Month() == t.Month()
		sameDay   = now.Day() == day
	)
	return sameYear && sameMonth && sameDay
}

func parseTime(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}

	t, err := time.Parse("2006/1", *monthYear)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return t, nil
}
