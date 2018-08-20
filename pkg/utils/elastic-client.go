package utils

import (
	"fmt"
	"github.com/olivere/elastic"
)

// NewResult creates a new instance of Result
func NewElasticClient(url string) (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil {
		return nil, fmt.Errorf("could not get ElasticSearch client: %s", err)
	}

	return client, err
}
