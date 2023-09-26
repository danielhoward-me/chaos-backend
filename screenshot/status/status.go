package status

type Status int

// Needs to match frontend in src/constants.ts
const (
	Generated Status = iota
	Failed
	Generating
	InQueue
	NotInQueue
)
