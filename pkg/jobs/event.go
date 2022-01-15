package jobs

const (
	AddJob = iota
	DelJob
	ReloadJobList
)

type Event struct {
	Action int
	Value  string
	Job    Job
}
