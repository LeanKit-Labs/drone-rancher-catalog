package types

// Build specifies run-time build information from Drone
type Build struct {
	Number    int
	Workspace string
	Commit    string
	Branch    string
}
