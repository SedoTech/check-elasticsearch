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
	checkStringQueryCmd struct {
		out      io.Writer
		Client   elastic.Client
		Query    string
		Warning  string
		Critical string
		Index    string
		Debug    bool
	}
)

func newCheckStringQueryCmd(out io.Writer) *cobra.Command {
	c := &checkStringQueryCmd{out: out}

	cmd := &cobra.Command{
		Use:          "stringQuery",
		Short:        "check if an ElasticSearch string query result meets the thresholds",
		SilenceUsage: false,
		Args:         NameArgs(),
		PreRun: func(cmd *cobra.Command, args []string) {
			c.Query = args[0]
			client, err := utils.NewElasticClient(utils.Elk01)
			if err != nil {
				icinga.NewResult("NewElasticClient", icinga.ServiceStatusUnknown, err.Error()).Exit()
			}
			c.Client = *client
		},
		Run: func(cmd *cobra.Command, args []string) {
			c.run()
		},
	}

	cmd.Flags().StringVarP(&c.Critical, "critical", "c", "1:", "critical threshold for minimum amount of search results")
	cmd.Flags().StringVarP(&c.Warning, "warning", "w", "2:", "warning threshold for minimum amount of search results")
	cmd.Flags().StringVarP(&c.Index, "index", "i", "*", "the index to search in")
	cmd.Flags().BoolVarP(&c.Debug, "debug", "d", false, "switch debug mode on/off")

	return cmd
}

func (c *checkStringQueryCmd) run() {
	checkcheckStringQuery := queries.NewCheckStringQuery(c.Client, c.Query)
	results := checkcheckStringQuery.CheckStringQueryString(queries.CheckStringQueryOptions{
		ThresholdWarning:  c.Warning,
		ThresholdCritical: c.Critical,
		Index:             c.Index,
		Debug:             c.Debug,
	})
	results.Exit()
}
