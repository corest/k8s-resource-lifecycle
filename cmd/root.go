package cmd

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	name        = "k8s-resource-lifecycle"
	description = "CLI tool to generate Kubernetes resource history based on audit logs."
)

type Config struct {
	Logger *log.Logger
	Stderr io.Writer
	Stdout io.Writer

	//
	GitCommit string
	Source    string
}

func New(config Config) (*cobra.Command, error) {
	if config.Stderr == nil {
		config.Stderr = os.Stderr
	}
	if config.Stdout == nil {
		config.Stdout = os.Stdout
	}

	f := &flag{}

	r := &runner{
		flag:   f,
		logger: config.Logger,
		stderr: config.Stderr,
		stdout: config.Stdout,
	}

	c := &cobra.Command{
		Use:               name,
		Short:             description,
		Long:              description,
		PersistentPreRunE: r.PersistentPreRun,
		RunE:              r.Run,
		SilenceUsage:      true,
	}

	f.Init(c)

	return c, nil
}
