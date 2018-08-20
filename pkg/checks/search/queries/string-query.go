package queries

import (
	"fmt"
	"context"
	"github.com/benkeil/icinga-checks-library"
	"github.com/olivere/elastic"
	"encoding/json"
)

type (
	// StringQueryCheck interface to check a query string
	StringQueryCheck interface {
		StringQuery(StringQueryOptions) icinga.Result
	}

	stringQueryImpl struct {
		Client elastic.Client
		Query  string
	}
)

// NewStringQuery creates a new instance of StringQueryCheck
func NewStringQuery(client elastic.Client, query string) StringQueryCheck {
	return &stringQueryImpl{Client: client, Query: query}
}

// StringQueryOptions contains options needed to run StringQueryCheck check
type StringQueryOptions struct {
	Query             string
	ThresholdWarning  string
	ThresholdCritical string
	Index             string
	Cache             bool
	Verbose           int
}

// StringQuery checks the result of the given Lucene query against the threshold values
func (c *stringQueryImpl) StringQuery(options StringQueryOptions) icinga.Result {
	name := "StringQueryCheck.StringQuery"

	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}

	query := elastic.NewQueryStringQuery(c.Query)
	query.TimeZone("Europe/Berlin")

	searchResult, err := c.Client.Search().
		Index(options.Index). // Use specific elasticsearch index
		RequestCache(options.Cache). // Whether or not to use resluts from cache
		Query(query).
		From(0).Size(0). // Start the search from specific index. Return specific number of search hits.
		Pretty(true). // Pretty print JSON output
		Do(context.Background()) // Execute the search and return a SearchResult, using the Background context which enables the request to carry data through the process
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't query ElasticSearch: %v", err))
	}

	totalHits := searchResult.Hits.TotalHits
	status := statusCheck.Check(float64(totalHits))
	message := fmt.Sprintf("Search produced %v hit(s) - (Query took %d ms)", totalHits, searchResult.TookInMillis)

	if options.Verbose > 0 {
		src, err := query.Source()
		if err == nil {
			data, err := json.Marshal(src)
			if err == nil {
				fmt.Printf("%s: %v\n", name, string(data))
			}
		}

		cacheOnOff := "off"
		if options.Cache {
			cacheOnOff = "on"
		}
		fmt.Printf("Cache is switched %v\n", cacheOnOff)
		fmt.Printf("Query took %v ms\n", searchResult.TookInMillis)
		if totalHits > 0 {
			fmt.Printf("Description %v\n", searchResult.Hits.Hits[0].Explanation.Description)
		}
	}

	return icinga.NewResult(name, status, message)
}
