package main

import (
	"errors"
	"os"
	"github.com/spf13/cobra"
)

var (
	globalUsage = `This program queries ElasticSearch`
	version     string
)

func main() {
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
		newCheckStringQueryCmd(out),
	)

	return cmd
}

// NameArgs returns an error if there are not exactly 1 arg containing the resource name.
func NameArgs() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("query is required")
		}
		return nil
	}
}
