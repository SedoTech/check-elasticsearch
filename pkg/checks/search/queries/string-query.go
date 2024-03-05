package queries

import (
	"encoding/json"
	"fmt"
	"github.com/benkeil/icinga-checks-library"
	"github.com/elastic/go-elasticsearch/v7"
)

type (
	// StringQueryCheck interface to check a query string
	StringQueryCheck interface {
		StringQuery(StringQueryOptions) icinga.Result
	}

	stringQueryImpl struct {
		Client *elasticsearch.Client
		Query  string
	}
)

// NewStringQuery creates a new instance of StringQueryCheck
func NewStringQuery(client *elasticsearch.Client, query string) StringQueryCheck {
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
	client := c.Client
	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}
	searchResult, err := client.Search(
		client.Search.WithQuery(options.Query),
		client.Search.WithIndex(options.Index))
	if err != nil {
		return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
	}

	var result map[string]interface{}
	json.NewDecoder(searchResult.Body).Decode(&result)
	hitsMap, ok := result["hits"].(map[string]interface{})
	if !ok {
		fmt.Println("Invalid type for 'hits'")
	}
	totalHits, ok := hitsMap["total"].(float64)
	status := statusCheck.Check(float64(totalHits))
	//message := fmt.Sprintf("Search produced %v hit(s) - (Query took %d ms)", totalHits, searchResult.TookInMillis)
	message := fmt.Sprintf("Search produced %v hit(s) - (Query took %d ms)", int(totalHits), 0)
	/*
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
			//fmt.Printf("Query took %v ms\n", searchResult.TookInMillis)
			//if totalHits > 0 {
			//	fmt.Printf("Description %v\n", searchResult.Hits.Hits[0].Explanation.Description)
			//}
		}
	*/
	return icinga.NewResult(name, status, message)
}
