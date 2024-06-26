package cmd

import (
	"errors"
	"github.com/spf13/cobra"
	"os"
)

var (
	globalUsage = `This program queries ElasticSearch`
	version     string
)

func Execute() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:          "check-elasticsearch",
		Short:        "check-elasticsearch checks if an ElasticSearch query meets the thresholds",
		Long:         globalUsage,
		Version:      version,
		SilenceUsage: true,
	}

	cmd.PersistentFlags().Parse(args)

	out := cmd.OutOrStdout()

	cmd.AddCommand(
		// check commands
		newStringQueryCmd(out),
	)

	return cmd
}

// NameArgs returns an error if there are not exactly 1 arg containing the resource name.
func NameArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("URL is required")
		}
		return nil
	}
}
