package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type cronExpr struct {
	minute cronField
	hour   cronField
	dom    cronField
	month  cronField
	dow    cronField
}

type cronField struct {
	any    bool
	values map[int]bool
}

func parseCron(expr string) (cronExpr, error) {
	parts := strings.Fields(strings.TrimSpace(expr))
	if len(parts) != 5 {
		return cronExpr{}, fmt.Errorf("cron must have 5 fields: minute hour day-of-month month day-of-week")
	}
	minute, err := parseCronField(parts[0], 0, 59, false)
	if err != nil {
		return cronExpr{}, fmt.Errorf("minute: %w", err)
	}
	hour, err := parseCronField(parts[1], 0, 23, false)
	if err != nil {
		return cronExpr{}, fmt.Errorf("hour: %w", err)
	}
	dom, err := parseCronField(parts[2], 1, 31, false)
	if err != nil {
		return cronExpr{}, fmt.Errorf("day-of-month: %w", err)
	}
	month, err := parseCronField(parts[3], 1, 12, false)
	if err != nil {
		return cronExpr{}, fmt.Errorf("month: %w", err)
	}
	dow, err := parseCronField(parts[4], 0, 7, true)
	if err != nil {
		return cronExpr{}, fmt.Errorf("day-of-week: %w", err)
	}
	return cronExpr{minute: minute, hour: hour, dom: dom, month: month, dow: dow}, nil
}

func parseCronField(raw string, min, max int, dow bool) (cronField, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return cronField{}, fmt.Errorf("empty field")
	}
	out := cronField{values: make(map[int]bool)}
	if raw == "*" {
		out.any = true
		return out, nil
	}

	for _, item := range strings.Split(raw, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			return cronField{}, fmt.Errorf("empty list item")
		}
		base, step := item, 1
		if strings.Contains(item, "/") {
			parts := strings.Split(item, "/")
			if len(parts) != 2 {
				return cronField{}, fmt.Errorf("invalid step %q", item)
			}
			base = parts[0]
			n, err := strconv.Atoi(parts[1])
			if err != nil || n < 1 {
				return cronField{}, fmt.Errorf("invalid step %q", item)
			}
			step = n
		}

		start, end := min, max
		switch {
		case base == "*" || base == "":
		case strings.Contains(base, "-"):
			r := strings.Split(base, "-")
			if len(r) != 2 {
				return cronField{}, fmt.Errorf("invalid range %q", item)
			}
			a, err1 := strconv.Atoi(r[0])
			b, err2 := strconv.Atoi(r[1])
			if err1 != nil || err2 != nil || a < min || b > max || a > b {
				return cronField{}, fmt.Errorf("invalid range %q", item)
			}
			start, end = a, b
		default:
			n, err := strconv.Atoi(base)
			if err != nil || n < min || n > max {
				return cronField{}, fmt.Errorf("invalid value %q", item)
			}
			start, end = n, n
		}

		for n := start; n <= end; n += step {
			if dow && n == 7 {
				out.values[0] = true
			} else {
				out.values[n] = true
			}
		}
	}
	return out, nil
}

func (f cronField) matches(v int) bool {
	if f.any {
		return true
	}
	return f.values[v]
}

func (c cronExpr) matches(t time.Time) bool {
	if !c.minute.matches(t.Minute()) || !c.hour.matches(t.Hour()) || !c.month.matches(int(t.Month())) {
		return false
	}
	domMatch := c.dom.matches(t.Day())
	dowMatch := c.dow.matches(int(t.Weekday()))
	if c.dom.any && c.dow.any {
		return true
	}
	if c.dom.any {
		return dowMatch
	}
	if c.dow.any {
		return domMatch
	}
	// Vixie-cron semantics: when both DOM and DOW are restricted, either may match.
	return domMatch || dowMatch
}

func nextCronTime(expr, timezone string, after time.Time) (time.Time, error) {
	c, err := parseCron(expr)
	if err != nil {
		return time.Time{}, err
	}
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	t := after.In(loc).Truncate(time.Minute).Add(time.Minute)
	deadline := t.AddDate(1, 0, 1)
	for !t.After(deadline) {
		if c.matches(t) {
			return t, nil
		}
		t = t.Add(time.Minute)
	}
	return time.Time{}, fmt.Errorf("no cron occurrence found within one year")
}
