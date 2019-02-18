package main

import (
	"math"
	"time"

	"github.com/faiface/pixel"
)

// powHalf calculates 0.5^(i)
func powHalf(i int) float64 {
	return math.Pow(2, float64(-i))
}

// pow2 calculates 2^i
func pow2(i int) float64 {
	return math.Pow(2, float64(i))
}

// clock is a convenience struct to contain all
// data needed to represent the clock
type clock struct {
	faces []*pixel.Sprite
	positions
}

// pos contains the positional data for a sprite.
type pos struct {
	position pixel.Vec
	angle    float64
}

// positions is a slice of the positional data for sprites.
type positions []pos

// calculate uses the current time to determine the position and rotation
// for each of the len(p) sprites.
func (p positions) calculate(now time.Time, radius float64, cfg Configuration) {

	for i := Year; i <= Second; i++ {
		angle := 0.0

		switch i {
		case Year:
			angle = float64(now.YearDay()-1) / 365.0

		case Month:
			// uses idiosyncrasy of time package to figure out the number
			// of days in the current month.
			// see: https://yourbasic.org/golang/last-day-month-date/
			daysInMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Day()
			angle = float64(now.Day()-1) / float64(daysInMonth)

		case Day:
			angle = float64(now.Hour()) / 24.0 // Hour() returns [0,23]

		case Hour:
			angle = float64(now.Minute()) / 60.0 // Minute() [0,59]

		case Minute:
			angle = float64(now.Second()) / 60.0 // Second() [0,59]

		case Second:
			angle = float64(now.Nanosecond()) / 1e9
		}

		p[i].angle = angle * TwoPi
		p[i].position = p.locForIndex(i, radius, cfg)
	}

}

// locForIndex calculates the position of the component at i given the initial
// radius.
//    loc = prev_loc + radius * 2^(-i) * Vec2(sincos(prev_angle))
// except when i==0, loc = (0,0)
func (p positions) locForIndex(i int, radius float64, cfg Configuration) pixel.Vec {
	if i == 0 {
		return pixel.ZV
	}

	angle := 0.0
	switch cfg.RotationMode {
	case DoNotRotate:
		// no nothing
	case AngleRelativeToZero:
		angle = p[i].angle
	case AngleIsParentAngle:
		angle = p[i-1].angle
	case AngleRelativeToParent:
		angle = p[i].angle + p[i-i].angle
	}

	// use sin for x and cos for y here to "rotate" everything by 90 degrees
	// making 0 be at "12 o'clock"
	dirUnitCircle := pixel.V(math.Sincos(-1 * float64(cfg.RotationDirection) * angle))
	return p[i-1].position.Add(dirUnitCircle.Scaled(radius * powHalf(i)))
}
