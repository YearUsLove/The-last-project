package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	nowTime, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, "invalid now date", http.StatusBadRequest)
		return
	}

	next, err := NextDate(nowTime, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(next))
}

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	date, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid format")
	}

	rule := strings.ToLower(strings.TrimSpace(parts[0]))

	switch rule {
	case "d":
		return handleDaysRule(date, now, parts)
	case "y":
		return handleYearsRule(date, now)
	case "w":
		return handleWeeksRule(date, now, parts)
	case "m":
		return handleMonthsRule(date, now, parts)
	default:
		return "", fmt.Errorf("invalid format")
	}
}

func handleDaysRule(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format")
	}
	n, err := strconv.Atoi(strings.ReplaceAll(parts[1], "+", ""))
	if err != nil || n < 1 || n > 400 {
		return "", fmt.Errorf("invalid days number")
	}

	date = date.AddDate(0, 0, n)

	for !date.After(now) {
		date = date.AddDate(0, 0, n)
	}

	return date.Format(dateFormat), nil
}

func handleYearsRule(date, now time.Time) (string, error) {
	original := date

	for !date.After(now) {
		year := date.Year() + 1
		if date.Month() == time.February && date.Day() == 29 && !isLeapYear(year) {
			date = time.Date(year, time.March, 1, 0, 0, 0, 0, time.UTC)
		} else {
			date = date.AddDate(1, 0, 0)
		}
	}

	if original.After(now) && date.Equal(original) {
		year := date.Year() + 1
		if date.Month() == time.February && date.Day() == 29 && !isLeapYear(year) {
			date = time.Date(year, time.March, 1, 0, 0, 0, 0, time.UTC)
		} else {
			date = date.AddDate(1, 0, 0)
		}
	}

	return date.Format(dateFormat), nil
}

func isLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}

func handleWeeksRule(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format")
	}
	days := strings.Split(parts[1], ",")
	weekdays := make([]int, 0, len(days))
	for _, d := range days {
		v, err := strconv.Atoi(strings.TrimSpace(d))
		if err != nil || v < 1 || v > 7 {
			return "", fmt.Errorf("invalid weekday")
		}
		if v == 7 {
			weekdays = append(weekdays, 0)
		} else {
			weekdays = append(weekdays, v)
		}
	}

	for {
		candidate := time.Time{}
		for _, wd := range weekdays {
			diff := (wd - int(date.Weekday()) + 7) % 7
			if diff == 0 {
				diff = 7
			}
			next := date.AddDate(0, 0, diff)
			if next.After(now) && (candidate.IsZero() || next.Before(candidate)) {
				candidate = next
			}
		}
		if !candidate.IsZero() {
			return candidate.Format(dateFormat), nil
		}
		date = date.AddDate(0, 0, 7)
	}
}

func handleMonthsRule(date, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("invalid format")
	}

	days := []int{}
	for _, d := range strings.Split(parts[1], ",") {
		v, err := strconv.Atoi(strings.TrimSpace(d))
		if err != nil || !(v >= 1 && v <= 31 || v == -1 || v == -2) {
			return "", fmt.Errorf("invalid month day")
		}
		days = append(days, v)
	}

	var months []time.Month
	if len(parts) == 3 {
		for _, m := range strings.Split(parts[2], ",") {
			v, err := strconv.Atoi(strings.TrimSpace(m))
			if err != nil || v < 1 || v > 12 {
				return "", fmt.Errorf("invalid month")
			}
			months = append(months, time.Month(v))
		}
	}

	for {
		candidate := time.Time{}
		for i := 0; i < 24; i++ {
			year := date.Year() + (int(date.Month())-1+i)/12
			month := time.Month((int(date.Month())-1+i)%12 + 1)

			if len(months) > 0 {
				found := false
				for _, allowed := range months {
					if month == allowed {
						found = true
						break
					}
				}
				if !found {
					continue
				}
			}

			for _, day := range days {
				var realDay int
				if day > 0 {
					if day > daysInMonth(year, month) {
						continue
					}
					realDay = day
				} else {
					last := daysInMonth(year, month)
					realDay = last + day + 1
					if realDay < 1 {
						continue
					}
				}

				next := time.Date(year, month, realDay, 0, 0, 0, 0, time.UTC)
				if next.After(now) && next.After(date) {
					if candidate.IsZero() || next.Before(candidate) {
						candidate = next
					}
				}
			}
		}
		if !candidate.IsZero() {
			return candidate.Format(dateFormat), nil
		}
		date = date.AddDate(1, 0, 0)
	}
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
