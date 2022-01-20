package job

type action int

const (
	Add action = iota
	Del
	Reload
)

type Event struct {
	Action action
	Value  string
	Job    *Job
}
