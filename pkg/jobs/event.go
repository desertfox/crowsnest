package jobs

type action int

const (
	AddJob action = iota
	DelJob
	ReloadJobList
)

type Event struct {
	Action action
	Value  string
	Job    Job
}
