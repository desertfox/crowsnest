package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/desertfox/crowsnest/api"
	"github.com/desertfox/crowsnest/pkg/crows"
	"github.com/spf13/cobra"
)

var (
	rootCmd    = &cobra.Command{Use: "crowsnest"}
	configPath string
	c          crows.Config = crows.Config{}
)

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "config")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if configPath != "" {
		if err := c.Load(configPath); err != nil {
			log.Fatal(err)
		}
	}

	n := c.BuildNest()

	n.Start()

	api.New(n).Start()
}
