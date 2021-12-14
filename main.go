package main

import (
	"fmt"
	opsys "os"
	"time"

	"github.com/go-co-op/gocron"
)

const configPathDefault string = "/etc/crowsnest/config.yaml"

var c config = loadConfig(configPathDefault)

func init() {
	if opsys.Getenv("CROWSNEST_USERNAME") == "" || opsys.Getenv("CROWSNEST_PASSWORD") == "" {
		fmt.Println("Missing CROWSNEST_USERNAME or CROWSNEST_PASSWORD ENV variable.")
		opsys.Exit(1)
	}

	c.InitSession(opsys.Getenv("CROWSNEST_USERNAME"), opsys.Getenv("CROWSNEST_PASSWORD"))
}

func main() {
	s := gocron.NewScheduler(time.UTC)

	for _, j := range c.jobs {
		s.Every(j.getFrequency()).Minutes().Do(j.getFunc(c))
	}

	s.StartBlocking()
}

func bailOut(err error) {
	fmt.Println(err.Error())
	opsys.Exit(1)
}
