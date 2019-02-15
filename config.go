package main

// Configuration options
type Configuration struct {
	TargetFPS         int
	RotationDirection int
	RotationMode      int
}

// DefaultConfig gets the default configuration. Duh.
func DefaultConfig() Configuration {
	return Configuration{
		TargetFPS:         60,
		RotationDirection: CW,
		RotationMode:      BaseOnParentCenterLine,
	}
}
