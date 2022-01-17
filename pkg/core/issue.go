package core

type Issue struct {
	Tag      string
	File     string
	Line     int
	Column   int
	Severity string
	Message  string
	Source   string
}
