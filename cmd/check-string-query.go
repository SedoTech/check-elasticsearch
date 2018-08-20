package main

import (
	"io"
	"github.com/benkeil/icinga-checks-library"
	"github.com/olivere/elastic"
	"github.com/sgnl04/check-elasticsearch/pkg/checks/search/queries"
	"github.com/sgnl04/check-elasticsearch/pkg/utils"
	"github.com/spf13/cobra"
)

type (
	stringQueryCmd struct {
		out      io.Writer
		Client   elastic.Client
		Url      string
		Query    string
		Warning  string
		Critical string
		Index    string
		Cache    bool
		Verbose  int
	}
)

func newStringQueryCmd(out io.Writer) *cobra.Command {
	c := &stringQueryCmd{out: out}

	cmd := &cobra.Command{
		Use:          "stringQuery [URL] [flags]",
		Short:        "check if an ElasticSearch string query result meets the thresholds",
		SilenceUsage: false,
		Args:         NameArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			c.Url = args[0]
			client, err := utils.NewElasticClient(c.Url)
			if err != nil {
				icinga.NewResult("NewElasticClient", icinga.ServiceStatusUnknown, err.Error()).Exit()
			}
			c.Client = *client
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}

	cmd.Flags().StringVarP(&c.Query, "query", "q", "", "the query to execute")
	cmd.Flags().StringVarP(&c.Critical, "critical", "c", "10:", "critical threshold for minimum amount of search results")
	cmd.Flags().StringVarP(&c.Warning, "warning", "w", "5:", "warning threshold for minimum amount of search results")
	cmd.Flags().StringVarP(&c.Index, "index", "i", "*", "the index to search in")
	cmd.Flags().BoolVarP(&c.Cache, "cache", "e", false, "switch using query cache on/off (default is cache off)")
	cmd.Flags().CountVarP(&c.Verbose, "verbose", "v", "enable verbose output")

	cmd.MarkFlagRequired("query")

	return cmd
}

func (c *stringQueryCmd) run() {
	stringQuery := queries.NewStringQuery(c.Client, c.Query)
	results := stringQuery.StringQuery(queries.StringQueryOptions{
		Query:             c.Query,
		ThresholdWarning:  c.Warning,
		ThresholdCritical: c.Critical,
		Index:             c.Index,
		Cache:             c.Cache,
		Verbose:           c.Verbose,
	})
	results.Exit()
}
