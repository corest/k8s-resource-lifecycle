package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/briandowns/spinner"
	"github.com/giantswarm/microerror"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/corest/k8s-resource-lifecycle/pkg/project"
)

type runner struct {
	flag   *flag
	logger *log.Logger
	stdout io.Writer
	stderr io.Writer
}

func (r *runner) PersistentPreRun(cmd *cobra.Command, args []string) error {
	fmt.Printf("Version = %#q\n", project.Version())
	fmt.Printf("Git SHA = %#q\n", project.GitSHA())
	fmt.Printf("Command = %#q\n", cmd.Name())
	fmt.Println()

	return nil
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return microerror.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)

	var auditLogFiles []string
	{
		s.Suffix = " Searching for audit log files"
		s.Start()
		auditLogFiles, err = auditLogFilesSearch(r.flag.AuditLogPath, r.flag.SearchPattern, r.flag.RecursiveSearch)
		if err != nil {
			return microerror.Mask(err)
		}
		s.Stop()
	}

	fmt.Println(auditLogFiles)
	return nil
}

func auditLogFilesSearch(path, searchPattern string, recursive bool) ([]string, error) {

	var result []string
	var err error

	searchRegexp, err := regexp.Compile(searchPattern)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	if recursive {
		err = filepath.Walk(path, func(path string, file os.FileInfo, err error) error {
			if err == nil && searchRegexp.MatchString(file.Name()) {
				result = append(result, path)
			}
			return nil
		})
		if err != nil {
			return nil, microerror.Mask(err)
		}

		return result, nil
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	for _, file := range files {
		if searchRegexp.MatchString(file.Name()) {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}

	return result, nil
}
