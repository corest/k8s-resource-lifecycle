package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/corest/k8s-resource-lifecycle/pkg/project"

	"github.com/corest/k8s-resource-lifecycle/cmd"
)

func main() {
	err := mainE(context.Background())
	if err != nil {
		panic(fmt.Sprintf("%#v\n", err))
	}
}

func mainE(ctx context.Context) error {
	var err error

	logger := log.New()

	var rootCommand *cobra.Command
	{
		c := cmd.Config{
			Logger: logger,

			GitCommit: project.GitSHA(),
			Source:    project.Source(),
		}

		rootCommand, err = cmd.New(c)
		if err != nil {
			logger.Errorf("failed to create command: %v", err)
			os.Exit(1)
		}
	}

	err = rootCommand.Execute()
	if err != nil {
		logger.Errorf("failed to execute command: %v", err)
		os.Exit(1)
	}

	return nil
}
