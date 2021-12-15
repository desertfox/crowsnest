package main

import (
	"fmt"
	"os"
	opsys "os"
	"time"

	"github.com/go-co-op/gocron"
)

var c config = loadConfig(os.Args[1])

func init() {
	if opsys.Getenv("CROWSNEST_USERNAME") == "" || opsys.Getenv("CROWSNEST_PASSWORD") == "" {
		fmt.Println("Missing CROWSNEST_USERNAME or CROWSNEST_PASSWORD ENV variable.")
		opsys.Exit(1)
	}

	c.InitSession(opsys.Getenv("CROWSNEST_USERNAME"), opsys.Getenv("CROWSNEST_PASSWORD"))
}

func main() {
	s := gocron.NewScheduler(time.UTC)

	for _, j := range c.Jobs {
		s.Every(j.getFrequency()).Minutes().Do(j.getFunc(c))
	}

	s.StartBlocking()
}

func bailOut(err error) {
	fmt.Println(err.Error())
	opsys.Exit(1)
}
