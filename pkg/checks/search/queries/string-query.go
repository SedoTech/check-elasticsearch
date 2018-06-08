package queries

import (
	"fmt"
	"context"
	"github.com/benkeil/icinga-checks-library"
	"github.com/olivere/elastic"
	"encoding/json"
)

type (
	// CheckStringQuery interface to check a query string
	CheckStringQuery interface {
		CheckStringQueryString(CheckStringQueryOptions) icinga.Result
	}

	checkStringQueryImpl struct {
		Client elastic.Client
		Query  string
	}
)

// NewCheckStringQuery creates a new instance of CheckStringQuery
func NewCheckStringQuery(client elastic.Client, query string) CheckStringQuery {
	return &checkStringQueryImpl{Client: client, Query: query}
}

// CheckAvailableAddressesOptions contains options needed to run CheckAvailableAddresses check
type CheckStringQueryOptions struct {
	ThresholdWarning  string
	ThresholdCritical string
	Index             string
	Debug             bool
}

// CheckAvailableAddresses checks if the deployment has a minimum of available replicas
func (c *checkStringQueryImpl) CheckStringQueryString(options CheckStringQueryOptions) icinga.Result {
	name := "Queries.StringQuery"

	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}

	query := elastic.NewQueryStringQuery(c.Query)
	query.TimeZone("Europe/Berlin")

	if options.Debug {
		src, err := query.Source()
		if err == nil {
			data, err := json.Marshal(src)
			if err == nil {
				fmt.Printf("NewQueryStringQuery: %v\n", string(data))
			}
		}
	}

	searchResult, err := c.Client.Search().
		Index(options.Index).
		RequestCache(false).
		Query(query).
		From(0).Size(0).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't query ElasticSearch: %v", err))
	}

	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	totalHits := searchResult.Hits.TotalHits
	status := statusCheck.Check(float64(totalHits))
	message := fmt.Sprintf("Search produced %v hit(s) - (Query took %d ms from index [%v])",
		totalHits, searchResult.TookInMillis, options.Index)

	return icinga.NewResult(name, status, message)
}
