package search

type Output struct {
	Verbose int           `yaml:"verbose"`
	URL     string        `yaml:"teamsurl"`
	Client  OutputService `yaml:"-"`
}
type OutputService interface {
	Send(string, string) error
}

func (o Output) IsVerbose() bool {
	return o.Verbose > 0
}

func (o Output) Send(url string, text string) {
	o.Client.Send(url, text)
}
