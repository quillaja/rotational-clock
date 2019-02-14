package main

import "math"

// TwoPi = 2*Pi
const TwoPi = math.Pi * 2

// These are used for array indices and various scale calculations.
const (
	Year = iota
	Month
	Day
	Hour
	Minute
	Second
)

// unused
// var daysInMonth = [12]int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
