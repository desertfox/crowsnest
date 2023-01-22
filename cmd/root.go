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
		log.Println("Loading config", configPath)
		err := c.Load(configPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Starting API")
	api.New(c.BuildNest()).Start()
}