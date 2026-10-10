package main

import (
	"testing"
	"time"
)

func TestCronMultipleHours(t *testing.T) {
	c, err := parseCron("0 6,11,16,21 * * *")
	if err != nil {
		t.Fatal(err)
	}
	loc := time.FixedZone("test", 8*3600)
	yes := []time.Time{
		time.Date(2026, 10, 1, 6, 0, 0, 0, loc),
		time.Date(2026, 10, 1, 11, 0, 0, 0, loc),
		time.Date(2026, 10, 1, 16, 0, 0, 0, loc),
		time.Date(2026, 10, 1, 21, 0, 0, 0, loc),
	}
	for _, tm := range yes {
		if !c.matches(tm) {
			t.Fatalf("expected match: %v", tm)
		}
	}
	if c.matches(time.Date(2026, 10, 1, 11, 1, 0, 0, loc)) {
		t.Fatal("unexpected minute match")
	}
}

func TestCronStepAndSundaySeven(t *testing.T) {
	c, err := parseCron("*/15 * * * 7")
	if err != nil {
		t.Fatal(err)
	}
	sunday := time.Date(2026, 10, 4, 12, 30, 0, 0, time.UTC)
	monday := time.Date(2026, 10, 5, 12, 30, 0, 0, time.UTC)
	if !c.matches(sunday) {
		t.Fatal("expected Sunday match")
	}
	if c.matches(monday) {
		t.Fatal("unexpected Monday match")
	}
}

func TestNextCronTime(t *testing.T) {
	after := time.Date(2026, 10, 1, 10, 59, 20, 0, time.UTC)
	next, err := nextCronTime("0 11 * * *", "UTC", after)
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2026, 10, 1, 11, 0, 0, 0, time.UTC)
	if !next.Equal(want) {
		t.Fatalf("got %v want %v", next, want)
	}
}

func TestNextCronTimeIncludesSecondSend(t *testing.T) {
	after := time.Date(2026, 10, 10, 6, 0, 20, 0, time.UTC)
	next, err := nextCronTime("0 6,11,16,21 * * *", "UTC", after)
	want := time.Date(2026, 10, 10, 6, 1, 0, 0, time.UTC)
	if err != nil || !next.Equal(want) {
		t.Fatalf("got %v, %v; want %v", next, err, want)
	}
}
