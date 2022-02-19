package job

const (
	Add int = iota
	Del
	Reload
)

type Event struct {
	Action int
	Value  string
	Job    *Job
}
