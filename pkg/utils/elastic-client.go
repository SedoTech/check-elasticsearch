package utils

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

// NewResult creates a new instance of Result
func NewElasticClient(url string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	client, err := elasticsearch.NewClient(cfg)

	if err != nil {
		return nil, fmt.Errorf("could not get ElasticSearch client: %s", err)
	}

	return client, err
}
