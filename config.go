package main

// Configuration options
type Configuration struct {
	TargetFPS         int
	RotationDirection int
	RotationMode      int
	RadiusOffset      int
	BackgroundColor   string
}

// DefaultConfig gets the default configuration. Duh.
func DefaultConfig() Configuration {
	return Configuration{
		TargetFPS:         60,
		RotationDirection: CCW,
		RotationMode:      AngleRelativeToZero,
		RadiusOffset:      -37,
		BackgroundColor:   "#CCC",
	}
}
