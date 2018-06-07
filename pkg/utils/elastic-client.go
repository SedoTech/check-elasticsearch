package main

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"github.com/olivere/elastic"
	"github.com/sgnl04/check-elasticsearch/pkg/structs"
)

// NewResult creates a new instance of Result
func NewElasticClient(url string) elastic.Client {
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil {
		// Handle error
		panic(err)
	}

	return client
}

func main() {
	// Create a client
	client, err := elastic.NewClient(elastic.SetURL("http://elk01.i.sedorz.net:9200"))

	if err != nil {
		// Handle error
		panic(err)
	}

	// Search with a term query
	//termQuery := elastic.NewTermQuery("message", "login")
	//query := elastic.NewSimpleQueryStringQuery(`message:("activemq")`)
	query := elastic.NewQueryStringQuery(`message:("ldap") AND level:<7 AND @timestamp:>now-90min`)
	query.TimeZone("Europe/Berlin")
	searchResult, err := client.Search().Explain(true).
		Index("integration.integration-marketplace.ecommerce.logs-6.2.2-2018.06.06").            // search in index "tweets"
		Query(query). // specify the query
		Sort("@timestamp", false).
		From(0).Size(1).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)
	//fmt.Printf("Query Profile: %s\n", searchResult.)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over iterating the hits, see below.
	var ttyp structs.EcommerceTweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(structs.EcommerceTweet); ok {
			fmt.Printf("%s | Tweet by [%s|%s (%d)]: %s\n", t.Timestamp, t.Application, t.Severity, t.Level, t.Message)
		}
	}
	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d tweets\n", searchResult.TotalHits())

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d tweets\n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			fmt.Printf("hit.Source: %s\n", hit.Source)
			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t structs.EcommerceTweet
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			//fmt.Printf("Tweet by [%s]: %s\n", t.User/*, t.Severity, t.Index*/, t.Message)
		}
	} else {
		// No hits
		fmt.Print("Found no tweets\n")
	}

	// Delete the index again
	//_, err = client.DeleteIndex("tweets").Do(context.Background())
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
}