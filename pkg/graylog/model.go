package graylog

type sessionService interface {
	GetHeader() string
	GetHost() string
}

type queryService interface {
	Execute(string) (int, error)
	BuildSearchURL() string
}

type Graylog struct {
	sessionService
	queryService
}

func New(sessionService sessionService, queryService queryService) Graylog {
	return Graylog{sessionService, queryService}
}

func (g Graylog) ExecuteSearch() (int, error) {
	return g.Execute(g.GetHeader())
}
