package job

import (
	"fmt"
	"log"
)

type Job struct {
	Name      string    `yaml:"name"`
	Condition Condition `yaml:"condition"`
	Output    Output    `yaml:"output"`
	Search    Search    `yaml:"search"`
}

func (j Job) Func(searchService SearchService, reportService OutputService) func() {
	return func() {
		j := j

		rawSearch, err := searchService.Execute()
		if err != nil {
			log.Fatal(err.Error())
		}

		j.Condition.Eval(rawSearch)

		if j.Output.IsVerbose() || j.Condition.IsAlert() {
			reportService.Send(
				j.Output.TeamsURL,
				j.buildText(searchService.BuildURL()),
			)
		}
	}
}

func (j Job) buildText(url string) string {
	return fmt.Sprintf("🔎 Name  : %s\n\r"+
		"⌚ Freq  : %d\n\r"+
		"📜 Status: %s\n\r"+
		"🧮 Count : %d\n\r"+
		"🔗 Link  : [GrayLog](%s)",
		j.Name,
		j.Search.Frequency,
		j.Condition.IsAlertText(),
		j.Condition.Count,
		url,
	)
}
