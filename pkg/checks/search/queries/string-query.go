package queries

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/SedoTech/icinga-checks-library"
	"github.com/elastic/go-elasticsearch/v7"
)

type (
	StringQueryCheck interface {
		StringQuery(StringQueryOptions) icinga.Result
	}

	stringQueryImpl struct {
		Client *elasticsearch.Client
		Query  string
	}
)

func NewStringQuery(client *elasticsearch.Client, query string) StringQueryCheck {
	return &stringQueryImpl{Client: client, Query: query}
}

type StringQueryOptions struct {
	Query             string
	ThresholdWarning  string
	ThresholdCritical string
	Index             string
	Cache             bool
	Verbose           int
}

func (c *stringQueryImpl) StringQuery(options StringQueryOptions) icinga.Result {
	name := "StringQueryCheck.StringQuery"
	client := c.Client
	statusCheck, err := icinga.NewStatusCheck(options.ThresholdWarning, options.ThresholdCritical)
	if err != nil {
		return icingaErrorMessage(name, err)
	}

	searchResult, err := client.Search(
		client.Search.WithQuery(options.Query),
		client.Search.WithIndex(options.Index),
		client.Search.WithRequestCache(options.Cache))
	if err != nil {
		return icingaErrorMessage(name, err)
	}

	totalHits, timeTook, err := extractHitsInfo(searchResult.Body)
	if err != nil {
		return icingaErrorMessage(name, err)
	}
	status := statusCheck.Check(float64(totalHits))
	message := fmt.Sprintf("Search produced %v hit(s) - (Query took %d ms)", int(totalHits), int(timeTook))

	// not sure what this tries to accomplish tbh.
	if options.Verbose > 0 {
		cacheOnOff := "off"
		if options.Cache {
			cacheOnOff = "on"
		}
		fmt.Printf("Cache is switched %v\n", cacheOnOff)
		fmt.Printf("Query took %v ms\n", timeTook)
	}
	return icinga.NewResult(name, status, message)
}

func extractHitsInfo(body io.ReadCloser) (float64, float64, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return 0, 0, err
	}

	hitsMap, ok := result["hits"].(map[string]interface{})
	if !ok {
		return 0, 0, fmt.Errorf("Invalid type for 'hits'")
	}

	totalHits, err := extractTotalHits(hitsMap)
	if err != nil {
		return 0, 0, err
	}

	timeTook, err := extractFloat64(result, "took")
	if err != nil {
		return 0, 0, err
	}

	return totalHits, timeTook, nil
}

func extractFloat64(data map[string]interface{}, key string) (float64, error) {
	value, ok := data[key].(float64)
	if !ok {
		return 0, fmt.Errorf("Invalid type for '%s'", key)
	}
	return value, nil
}

func extractTotalHits(data map[string]interface{}) (float64, error) {
	total := data["total"]
	if value, ok := total.(map[string]interface{}); ok {
		return extractFloat64(value, "value")
	}
	return extractFloat64(data, "total")
}

func icingaErrorMessage(name string, err error) icinga.Result {
	return icinga.NewResult(name, icinga.ServiceStatusUnknown, fmt.Sprintf("can't check status: %v", err))
}
