package handlers

import (
	"fmt"
	"time"
)

// ── Gregorian ─────────────────────────────────────────────────────────────────

func gregorianDate(t time.Time) string {
	return t.Format("Monday, 2 January 2006")
}

// ── Hijri (Islamic / Arabic) ──────────────────────────────────────────────────
// Uses the tabular Islamic calendar algorithm (Fliegel & Van Flandern).

var hijriMonths = [12]string{
	"مُحَرَّم", "صَفَر", "رَبِيع الأَوَّل", "رَبِيع الثَّانِي",
	"جُمَادَى الأُولَى", "جُمَادَى الآخِرَة", "رَجَب", "شَعْبَان",
	"رَمَضَان", "شَوَّال", "ذُو القَعْدَة", "ذُو الحِجَّة",
}

func hijriDate(t time.Time) string {
	// Convert Gregorian to Julian Day Number (JDN)
	y, m, d := t.Year(), int(t.Month()), t.Day()
	a := (14 - m) / 12
	yy := y + 4800 - a
	mm := m + 12*a - 3
	jdn := d + (153*mm+2)/5 + 365*yy + yy/4 - yy/100 + yy/400 - 32045

	// Convert JDN to Hijri
	l := jdn - 1948440 + 10632
	n := (l - 1) / 10631
	l = l - 10631*n + 354
	j := ((10985-l)/5316)*((50*l)/17719) + (l/5670)*((43*l)/15238)
	l = l - ((30-j)/15)*((17719*j)/50) - (j/16)*((15238*j)/43) + 29
	hm := (24 * l) / 709
	hd := l - (709*hm)/24
	hy := 30*n + j - 30

	return fmt.Sprintf("%d %s %d هـ", hd, hijriMonths[hm-1], hy)
}

// ── Bengali (Bangla Academy revised calendar) ─────────────────────────────────
// Month boundaries follow the Bangla Academy 1987 reform (fixed Gregorian dates).
// Falgun and Chaitra start one day earlier in Gregorian leap years.

var bengaliMonths = [12]string{
	"বৈশাখ", "জ্যৈষ্ঠ", "আষাঢ়", "শ্রাবণ",
	"ভাদ্র", "আশ্বিন", "কার্তিক", "অগ্রহায়ণ",
	"পৌষ", "মাঘ", "ফাল্গুন", "চৈত্র",
}

func isLeapYear(y int) bool {
	return (y%4 == 0 && y%100 != 0) || y%400 == 0
}

func bengaliDate(t time.Time) string {
	gy, gm, gd := t.Year(), int(t.Month()), t.Day()
	tDate := time.Date(gy, time.Month(gm), gd, 0, 0, 0, 0, time.UTC)

	// Bengali year: starts April 14, so dates before Apr 14 belong to previous Bengali year
	var bYear int
	if gm > 4 || (gm == 4 && gd >= 14) {
		bYear = gy - 593
	} else {
		bYear = gy - 594
	}

	g1 := bYear + 593 // Gregorian year containing Baishakh–Poush (Apr–Dec)
	g2 := bYear + 594 // Gregorian year containing Magh–Chaitra (Jan–Mar)

	// Falgun and Chaitra start one day earlier when g2 is a Gregorian leap year
	falgunDay, chaitraDay := 13, 14
	if isLeapYear(g2) {
		falgunDay, chaitraDay = 12, 13
	}

	// Month start dates in chronological order (Apr g1 → Mar g2)
	type monthBound struct {
		t   time.Time
		idx int // index into bengaliMonths
	}
	bounds := []monthBound{
		{time.Date(g1, 4, 14, 0, 0, 0, 0, time.UTC), 0},            // Baishakh
		{time.Date(g1, 5, 15, 0, 0, 0, 0, time.UTC), 1},            // Jyaistha
		{time.Date(g1, 6, 15, 0, 0, 0, 0, time.UTC), 2},            // Asharh
		{time.Date(g1, 7, 16, 0, 0, 0, 0, time.UTC), 3},            // Shraban
		{time.Date(g1, 8, 16, 0, 0, 0, 0, time.UTC), 4},            // Bhadra
		{time.Date(g1, 9, 16, 0, 0, 0, 0, time.UTC), 5},            // Ashwin
		{time.Date(g1, 10, 16, 0, 0, 0, 0, time.UTC), 6},           // Kartik
		{time.Date(g1, 11, 15, 0, 0, 0, 0, time.UTC), 7},           // Agrahayan
		{time.Date(g1, 12, 15, 0, 0, 0, 0, time.UTC), 8},           // Poush
		{time.Date(g2, 1, 14, 0, 0, 0, 0, time.UTC), 9},            // Magh
		{time.Date(g2, 2, falgunDay, 0, 0, 0, 0, time.UTC), 10},    // Falgun
		{time.Date(g2, 3, chaitraDay, 0, 0, 0, 0, time.UTC), 11},   // Chaitra
	}

	// Find the last boundary ≤ tDate (linear scan, list is in chronological order)
	idx := 0
	var start time.Time
	for _, b := range bounds {
		if !tDate.Before(b.t) {
			idx = b.idx
			start = b.t
		}
	}

	bDay := int(tDate.Sub(start).Hours()/24) + 1
	return fmt.Sprintf("%d %s %d বঙ্গাব্দ", bDay, bengaliMonths[idx], bYear)
}
