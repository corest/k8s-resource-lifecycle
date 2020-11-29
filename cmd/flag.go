package cmd

import (
	"regexp"

	"github.com/giantswarm/microerror"
	"github.com/spf13/cobra"
)

const (
	flagAuditLogPath    = "audit-log-path"
	flagRecursiveSearch = "recursive-search"
	flagSearchPattern   = "search-pattern"

	// resource flags
	flagResourceKind       = "resource-kind"
	flagResourceName       = "resource-name"
	flagResourceNamespace  = "resource-namespace"
	flagResourceAPIGroup   = "resource-apigroup"
	flagResourceAPIVersion = "resource-apiversion"
)

type flag struct {
	AuditLogPath    string
	RecursiveSearch bool
	SearchPattern   string

	// Resource options
	ResourceKind       string
	ResourceName       string
	ResourceNamespace  string
	ResourceAPIGroup   string
	ResourceAPIVersion string
}

func (f *flag) Init(cmd *cobra.Command) {

	cmd.Flags().StringVarP(&f.AuditLogPath, flagAuditLogPath, "p", ".", "Path to the directory, where audit log files are located.")
	cmd.Flags().BoolVarP(&f.RecursiveSearch, flagRecursiveSearch, "r", false, "Recursive lookup for audit log files.")
	cmd.Flags().StringVarP(&f.SearchPattern, flagSearchPattern, "s", ".*", "Search pattern used for audit log files lookup.")

	// validate resource
	cmd.Flags().StringVar(&f.ResourceKind, flagResourceKind, "", "Resource kind (e.g. `cluster`)")
	cmd.Flags().StringVar(&f.ResourceName, flagResourceName, "", "Resource name")
	cmd.Flags().StringVar(&f.ResourceNamespace, flagResourceNamespace, "", "Resource namespace")
	cmd.Flags().StringVar(&f.ResourceAPIGroup, flagResourceAPIGroup, "", "Resource apigroup")
	cmd.Flags().StringVar(&f.ResourceAPIVersion, flagResourceAPIVersion, "", "Resource apiversion")
}

func (f *flag) Validate() error {
	if f.AuditLogPath == "" {
		f.AuditLogPath = "."
	}

	_ = regexp.MustCompile(f.SearchPattern)

	if f.ResourceName == "" {
		return microerror.Maskf(invalidFlagsError, "--%s must not be empty", flagResourceName)
	}

	if f.ResourceNamespace == "" {
		return microerror.Maskf(invalidFlagsError, "--%s must not be empty", flagResourceNamespace)
	}

	if f.ResourceKind == "" {
		return microerror.Maskf(invalidFlagsError, "--%s must not be empty", flagResourceKind)
	}

	if f.ResourceAPIGroup == "" {
		return microerror.Maskf(invalidFlagsError, "--%s must not be empty", flagResourceAPIGroup)
	}

	if f.ResourceAPIVersion == "" {
		return microerror.Maskf(invalidFlagsError, "--%s must not be empty", flagResourceAPIVersion)
	}

	return nil
}
