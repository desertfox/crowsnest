package search

type Teams struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}
type Output struct {
	Verbose int           `yaml:"verbose"`
	Teams   Teams         `yaml:"teams"`
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

func (o Output) URL() string {
	return o.Teams.Url
}
