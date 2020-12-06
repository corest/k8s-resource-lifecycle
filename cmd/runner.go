package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/giantswarm/microerror"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apiserver/pkg/apis/audit"

	"github.com/corest/k8s-resource-lifecycle/pkg/metaresource"
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
		s.Prefix = "Searching for audit log files "
		s.Start()
		auditLogFiles, err = auditLogFilesSearch(r.flag.AuditLogPath, r.flag.SearchPattern, r.flag.RecursiveSearch)
		if err != nil {
			return microerror.Mask(err)
		}
		s.Stop()
	}

	fmt.Printf("Looking for resource events:\n\t- kind: %#q\n\t- name: %#q\n\t- apiGroup: %#q\n\n",
		r.flag.ResourceKind, r.flag.ResourceName, r.flag.ResourceAPIGroup)
	s.Prefix = "Processing log files "
	s.Start()

	// write constructor for this
	metaResource := metaresource.MetaResource{
		Kind:      r.flag.ResourceKind,
		Name:      r.flag.ResourceName,
		Namespace: r.flag.ResourceNamespace,
		APIGroup:  r.flag.ResourceAPIGroup,
	}

	var wg sync.WaitGroup
	storeCh := make(chan audit.Event)
	errCh := make(chan error)
	for _, f := range auditLogFiles {
		wg.Add(1)
		go metaResource.FindEvents(f, storeCh, errCh, &wg)
	}

	var errors []error
	go func() {
		for {
			select {
			case event := <-storeCh:
				metaResource.StoreEvent(event)
			case err := <-errCh:
				errors = append(errors, err)
			}
		}
	}()

	wg.Wait()
	s.Stop()

	for _, err := range errors {
		fmt.Println("error")
		return err
	}

	if len(metaResource.Events) == 0 {
		fmt.Println("No events found")
		return nil
	}

	metaResource.SortEvents()

	for _, event := range metaResource.Events {
		fmt.Printf("Found %#q resource event (%s)\n", event.Verb, event.StageTimestamp.String())
	}

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
