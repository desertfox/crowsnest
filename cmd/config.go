package cmd

import (
	"log"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/spf13/cobra"
)

func runWithConfig() *cobra.Command {
	configCmd := &cobra.Command{
		Use:  "config",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p := args[0]

			c := crows.Config{}

			if err := c.Load(p); err != nil {
				log.Fatalf("unable to load config %s, %s", p, err.Error())
			}

			n := c.BuildNest()

			if err := n.Start(); err != nil {
				log.Fatal(err)
			}

			api.New(n).Start()
		},
	}

	return configCmd
}
