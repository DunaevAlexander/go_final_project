package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	nowStr := r.FormValue("now")
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid now: "+err.Error())
			return
		}
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	next, err := NextDate(now, date, repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", err
	}
	if repeat == "" {
		return "", fmt.Errorf("empty repeat")
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty repeat")
	}

	switch parts[0] {
	case "d":
		return nextDaily(now, start, parts)
	case "y":
		return nextYearly(now, start, parts)
	case "w":
		return nextWeekly(now, start, parts)
	case "m":
		return nextMonthly(now, start, parts)
	}
	return "", fmt.Errorf("unsupported repeat format: %q", repeat)
}

func afterNow(date, now time.Time) bool {
	return date.After(now)
}

func nextDaily(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid d rule")
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil || n <= 0 || n > 400 {
		return "", fmt.Errorf("invalid d interval: %q", parts[1])
	}
	start = start.AddDate(0, 0, n)
	for !afterNow(start, now) {
		start = start.AddDate(0, 0, n)
	}
	return start.Format(dateFormat), nil
}

func nextYearly(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", fmt.Errorf("invalid y rule")
	}
	start = start.AddDate(1, 0, 0)
	for !afterNow(start, now) {
		start = start.AddDate(1, 0, 0)
	}
	return start.Format(dateFormat), nil
}

func nextWeekly(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid w rule")
	}
	allowed := make(map[int]bool)
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || d < 1 || d > 7 {
			return "", fmt.Errorf("invalid weekday: %q", s)
		}
		allowed[d] = true
	}
	date := start
	for {
		wd := int(date.Weekday())
		if wd == 0 {
			wd = 7
		}
		if afterNow(date, now) && allowed[wd] {
			return date.Format(dateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

func nextMonthly(now, start time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid m rule")
	}

	days := make(map[int]bool)
	for _, s := range strings.Split(parts[1], ",") {
		d, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil || d == 0 || d > 31 || d < -2 {
			return "", fmt.Errorf("invalid day of month: %q", s)
		}
		days[d] = true
	}

	months := make(map[int]bool)
	if len(parts) == 3 {
		for _, s := range strings.Split(parts[2], ",") {
			m, err := strconv.Atoi(strings.TrimSpace(s))
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("invalid month: %q", s)
			}
			months[m] = true
		}
	}

	date := start
	limit := now.AddDate(5, 0, 0)
	for date.Before(limit) {
		if afterNow(date, now) &&
			(len(months) == 0 || months[int(date.Month())]) &&
			matchDay(date, days) {
			return date.Format(dateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}
	return "", fmt.Errorf("no matching date found for rule %q", strings.Join(parts, " "))
}

func matchDay(date time.Time, days map[int]bool) bool {
	day := date.Day()
	lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	for d := range days {
		switch {
		case d > 0:
			if d == day {
				return true
			}
		case d == -1:
			if day == lastDay {
				return true
			}
		case d == -2:
			if day == lastDay-1 {
				return true
			}
		}
	}
	return false
}
