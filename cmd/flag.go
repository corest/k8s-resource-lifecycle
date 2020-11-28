package cmd

import (
	"regexp"

	"github.com/spf13/cobra"
)

const (
	flagAuditLogPath    = "audit-log-path"
	flagRecursiveSearch = "recursive-search"
	flagSearchPattern   = "search-pattern"
)

type flag struct {
	AuditLogPath    string
	RecursiveSearch bool
	SearchPattern   string
}

func (f *flag) Init(cmd *cobra.Command) {

	cmd.Flags().StringVarP(&f.AuditLogPath, flagAuditLogPath, "p", ".", "Path to the directory, where audit log files are located.")
	cmd.Flags().BoolVarP(&f.RecursiveSearch, flagRecursiveSearch, "r", false, "Recursive lookup for audit log files.")
	cmd.Flags().StringVarP(&f.SearchPattern, flagSearchPattern, "s", ".*", "Search pattern used for audit log files lookup.")

}

func (f *flag) Validate() error {
	if f.AuditLogPath == "" {
		f.AuditLogPath = "."
	}

	_ = regexp.MustCompile(f.SearchPattern)

	return nil
}
